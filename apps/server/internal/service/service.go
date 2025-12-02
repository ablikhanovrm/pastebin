package service

import "github.com/ablikhanovrm/pastebin/internal/repository"

type Services struct {
}

func (*Services) NewService(repo *repository.Repository) *Services {
	return &Services{}
}
