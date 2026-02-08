package paste

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
)

func mapPasteFromDB(p dbgen.Paste) *paste.Paste {
	return &paste.Paste{
		Uuid:       p.Uuid,
		Id:         p.ID,
		UserId:     p.UserID,
		Title:      p.Title,
		S3Key:      p.S3Key,
		Syntax:     paste.Syntax(p.Syntax),
		Visibility: paste.Visibility(p.Visibility),
		ViewsCount: p.ViewsCount,
		MaxViews:   p.MaxViews,
		ExpiresAt:  &p.ExpireAt.Time,
		CreatedAt:  p.CreatedAt.Time,
		UpdatedAt:  p.UpdatedAt.Time,
	}
}
