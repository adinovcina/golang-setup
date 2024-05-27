package redisstore

import (
	"context"
	"time"

	"github.com/adinovcina/golang-setup/tools/utils"
	"github.com/twinj/uuid"
)

// SetSession - sets user session. Expects userID, sessionID and Client name.
func (s *RedisStore) SetSession(ctx context.Context, uid uuid.UUID, sid, v string, redisTokenTTL time.Duration) error {
	// Since refresh token data is stored in claim as well, we will keep it in Redis 10 times longer than
	// actual Token expiration time
	return s.redis.Set(ctx, utils.FormatSessionKey(uid, sid), v, redisTokenTTL).Err()
}

// GetSession - gets user session. Expects userID and sessionID.
func (s *RedisStore) GetSession(ctx context.Context, uid uuid.UUID, sid string) (string, error) {
	value, err := s.redis.Get(ctx, utils.FormatSessionKey(uid, sid)).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

// DelSession - del user session. Expects userID and sessionID.
func (s *RedisStore) DelSession(ctx context.Context, uid uuid.UUID, sid string) error {
	return s.redis.Del(ctx, utils.FormatSessionKey(uid, sid)).Err()
}

// DelSessionWithKey - del user session. Expects session key.
func (s *RedisStore) DelSessionWithKey(ctx context.Context, key string) error {
	return s.redis.Del(ctx, key).Err()
}
