package paste

import (
	"github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/rs/zerolog"
)

type PasteService struct {
	repo   paste.PasteRepository
	logger zerolog.Logger
}

func NewPasteService(repo paste.PasteRepository, logger zerolog.Logger) *PasteService {
	return &PasteService{repo: repo, logger: logger}
}
