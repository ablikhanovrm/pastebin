package paste

import (
	"github.com/ablikhanovrm/pastebin/internal/repository/paste"
)

type PasteService struct {
	repo paste.PasteRepository
}

func NewPasteService(repo paste.PasteRepository) *PasteService {
	return &PasteService{repo: repo}
}
