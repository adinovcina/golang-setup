package account

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adinovcina/golang-setup/api"
	"github.com/adinovcina/golang-setup/store"
	"github.com/adinovcina/golang-setup/tools/utils"
	"github.com/twinj/uuid"
)

// createToken will generate claims and sign in response which will go into header.
func (s *service) createToken(ctx context.Context, in *store.User) (token string, err error) {
	userData := &api.Data{
		UserID:     in.ID,
		Email:      in.Email,
		Active:     in.Active,
		Role:       in.Role,
		UserRoleID: in.RoleID,
	}

	// If Session is not created then notify clients but does not expose issue
	jwtClaim, err := s.createSessionData(ctx, userData, s.conf.Redis.TokenTTL)
	if err != nil {
		return "", err
	}

	// Generate JWT Token
	token, err = jwtClaim.CreateToken(s.conf.MFA.AccessTokenExpiration, s.conf.Redis.SecretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// createTemporaryToken will generate and persist short lasting token.
func (s *service) createRefreshToken(userID uuid.UUID) (string, error) {
	refreshToken := api.NewRefreshToken()

	// Persist temp token
	err := s.repo.AddLoginToken(userID, int64(s.conf.MFA.RefreshTokenExpiration.Minutes()), refreshToken, store.GetTokenTypes().RefreshToken)

	return refreshToken, err
}

// createTemporaryToken will generate and persist short lasting token.
func (s *service) createTemporaryToken(user *store.User) (string, error) {
	// Create temp token
	temporaryToken := utils.GenerateUniqueID() + utils.GenerateUniqueID()
	// Persist temp token
	err := s.repo.AddLoginToken(user.ID, int64(s.conf.MFA.TemporaryTokenExpiration.Minutes()), temporaryToken, store.GetTokenTypes().MFA)

	return temporaryToken, err
}

// CreateSessionData will create a session object and store in redis.
func (s *service) createSessionData(ctx context.Context, userData *api.Data, redisTokenTTL time.Duration) (claim *api.Claim, err error) {
	claim = new(api.Claim)

	claim.NewID()
	claim.UserID = userData.UserID

	userData.SessionKey = utils.FormatSessionKey(userData.UserID, claim.SessionID)

	userDataMarshaled, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	err = s.inMemRepo.SetSession(ctx, userData.UserID, claim.SessionID, string(userDataMarshaled), redisTokenTTL)
	if err != nil {
		return nil, err
	}

	return claim, nil
}
