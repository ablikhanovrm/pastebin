package paste

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type PasteRepository interface {
	Create(ctx context.Context, u *paste.Paste) (*paste.Paste, error)
	GetByID(ctx context.Context, id uuid.UUID) (*paste.Paste, error)
	GetPastes(ctx context.Context, userId int64) ([]*paste.Paste, error)
	GetMyPastes(ctx context.Context, userId int64) ([]*paste.Paste, error)
	Update(ctx context.Context, paste *paste.Paste) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SqlcPasteRepository struct {
	q   *dbgen.Queries
	log zerolog.Logger
}

func NewSqlcPasteRepository(db dbgen.DBTX, log zerolog.Logger) *SqlcPasteRepository {
	return &SqlcPasteRepository{q: dbgen.New(db), log: log}
}

func (r *SqlcPasteRepository) Create(ctx context.Context, u *paste.Paste) (*paste.Paste, error) {
	var expire pgtype.Timestamptz
	if u.ExpiresAt != nil {
		expire = pgtype.Timestamptz{
			Time:  *u.ExpiresAt,
			Valid: true,
		}
	} else {
		expire = pgtype.Timestamptz{
			Valid: false,
		}
	}

	row, err := r.q.CreatePaste(ctx, dbgen.CreatePasteParams{
		Uuid:       u.Uuid,
		UserID:     u.UserId,
		Title:      u.Title,
		S3Key:      u.S3Key,
		MaxViews:   u.MaxViews,
		Syntax:     string(u.Syntax),
		Visibility: string(u.Visibility),
		ExpireAt:   expire,
	})

	if err != nil {
		return nil, err
	}

	return mapPasteFromDB(row), nil
}

func (r *SqlcPasteRepository) GetByID(ctx context.Context, pasteUuid uuid.UUID, userId int64) (*paste.Paste, error) {
	row, err := r.q.GetPasteById(ctx, dbgen.GetPasteByIdParams{
		Uuid:   pasteUuid,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	return mapPasteFromDB(row), nil
}

func (r *SqlcPasteRepository) GetPastes(ctx context.Context, userId int64) ([]*paste.Paste, error) {
	rows, err := r.q.GetPastes(ctx, userId)

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))

	for _, row := range rows {
		pastes = append(pastes, mapPasteFromDB(row))
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) GetMyPastes(ctx context.Context, userId int64) ([]*paste.Paste, error) {
	rows, err := r.q.GetUserPastes(ctx, userId)

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))

	for _, row := range rows {
		pastes = append(pastes, mapPasteFromDB(row))
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) Update(ctx context.Context, paste *paste.Paste) error {
	var expire pgtype.Timestamptz
	if paste.ExpiresAt != nil {
		expire = pgtype.Timestamptz{
			Time:  *paste.ExpiresAt,
			Valid: true,
		}
	} else {
		expire = pgtype.Timestamptz{
			Valid: false,
		}
	}

	err := r.q.UpdatePaste(ctx, dbgen.UpdatePasteParams{
		Uuid:       paste.Uuid,
		Title:      paste.Title,
		Syntax:     string(paste.Syntax),
		Visibility: string(paste.Visibility),
		MaxViews:   paste.MaxViews,
		ExpireAt:   expire,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlcPasteRepository) Delete(ctx context.Context, pasteUuid uuid.UUID, userId int64) error {
	err := r.q.DeletePaste(ctx, dbgen.DeletePasteParams{
		Uuid:   pasteUuid,
		UserID: userId,
	})

	if err != nil {
		return err
	}

	return nil
}
