package account

import (
	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/store"
	mailjet "github.com/adinovcina/golang-setup/tools/mailjet"
)

type service struct {
	conf          *config.Config
	repo          store.Repository
	inMemRepo     store.InMemRepository
	mailjetClient *mailjet.Client
}

func newService(conf *config.Config,
	repo store.Repository,
	inMemRepo store.InMemRepository,
	mailjetClient *mailjet.Client,
) service {
	return service{
		conf,
		repo,
		inMemRepo,
		mailjetClient,
	}
}
