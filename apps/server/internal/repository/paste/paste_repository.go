package paste

import (
	"context"
	"errors"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type PasteRepository interface {
	Create(ctx context.Context, userId int64, u *paste.Paste) (*paste.Paste, error)
	GetByID(ctx context.Context, userId int64, id uuid.UUID) (*paste.Paste, error)
	GetManyByIDs(ctx context.Context, userId int64, ids []uuid.UUID) ([]*paste.Paste, error)
	GetPastesFirstPage(ctx context.Context, params GetPastesFirstPageParams) ([]*paste.Paste, error)
	GetPastesAfterCursor(ctx context.Context, params GetPastesAfterCursorParams) ([]*paste.Paste, error)
	GetMyPastesFirstPage(ctx context.Context, params GetMyPastesFirstPageParams) ([]*paste.Paste, error)
	GetMyPastesAfterCursor(ctx context.Context, params GetMyPastesAfterCursorParams) ([]*paste.Paste, error)
	Update(ctx context.Context, userId int64, paste *paste.Paste) error
	Delete(ctx context.Context, userId int64, id uuid.UUID) error
}

type SqlcPasteRepository struct {
	q   *dbgen.Queries
	log zerolog.Logger
}

func NewSqlcPasteRepository(db dbgen.DBTX, log zerolog.Logger) *SqlcPasteRepository {
	return &SqlcPasteRepository{q: dbgen.New(db), log: log}
}

func (r *SqlcPasteRepository) Create(ctx context.Context, userId int64, u *paste.Paste) (*paste.Paste, error) {
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
		UserID:     userId,
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

func (r *SqlcPasteRepository) GetByID(ctx context.Context, userId int64, pasteUuid uuid.UUID) (*paste.Paste, error) {
	row, err := r.q.GetPasteById(ctx, dbgen.GetPasteByIdParams{
		Uuid:   pasteUuid,
		UserID: userId,
	})

	if err != nil {
		return nil, err
	}

	return mapPasteFromDB(row), nil
}

func (r *SqlcPasteRepository) GetManyByIDs(ctx context.Context, userId int64, ids []uuid.UUID) ([]*paste.Paste, error) {
	rows, err := r.q.GetManyByIds(ctx, dbgen.GetManyByIdsParams{
		Ids:    ids,
		Userid: userId,
	})

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, len(rows))

	for i, row := range rows {
		pastes[i] = mapPasteFromDB(row)
	}
	return pastes, nil
}

func (r *SqlcPasteRepository) GetPastesFirstPage(ctx context.Context, params GetPastesFirstPageParams) ([]*paste.Paste, error) {
	rows, err := r.q.GetPastesFirstPage(ctx, dbgen.GetPastesFirstPageParams{
		UserID: params.UserId,
		Limit:  params.Limit,
	})

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))

	for i, row := range rows {
		pastes[i] = mapPasteFromDB(row)
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) GetPastesAfterCursor(ctx context.Context, params GetPastesAfterCursorParams) ([]*paste.Paste, error) {
	rows, err := r.q.GetPastesAfterCursor(ctx, dbgen.GetPastesAfterCursorParams{
		UserID:    params.UserId,
		Limit:     params.Limit,
		CreatedAt: pgtype.Timestamptz{Time: params.Cursor, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))
	for i, row := range rows {
		pastes[i] = mapPasteFromDB(row)
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) GetMyPastesFirstPage(ctx context.Context, params GetMyPastesFirstPageParams) ([]*paste.Paste, error) {
	rows, err := r.q.GetUserPastesFirstPage(ctx, dbgen.GetUserPastesFirstPageParams{
		UserID: params.UserId,
		Limit:  params.Limit,
	})

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))

	for i, row := range rows {
		pastes[i] = mapPasteFromDB(row)
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) GetMyPastesAfterCursor(ctx context.Context, params GetMyPastesAfterCursorParams) ([]*paste.Paste, error) {
	rows, err := r.q.GetUserPastesAfterCursor(ctx, dbgen.GetUserPastesAfterCursorParams{
		UserID: params.UserId,
		CreatedAt: pgtype.Timestamptz{
			Time:  params.Cursor,
			Valid: true,
		},
		Limit: params.Limit,
	})

	if err != nil {
		return nil, err
	}

	pastes := make([]*paste.Paste, 0, len(rows))

	for _, row := range rows {
		pastes = append(pastes, mapPasteFromDB(row))
	}

	return pastes, nil
}

func (r *SqlcPasteRepository) Update(ctx context.Context, userId int64, opts *paste.Paste) (*paste.Paste, error) {
	var expire pgtype.Timestamptz
	if opts.ExpiresAt != nil {
		expire = pgtype.Timestamptz{
			Time:  *opts.ExpiresAt,
			Valid: true,
		}
	} else {
		expire = pgtype.Timestamptz{
			Valid: false,
		}
	}

	update, err := r.q.UpdatePaste(ctx, dbgen.UpdatePasteParams{
		Uuid:       opts.Uuid,
		UserID:     userId,
		Title:      opts.Title,
		Syntax:     string(opts.Syntax),
		Visibility: string(opts.Visibility),
		MaxViews:   opts.MaxViews,
		ExpireAt:   expire,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, paste.ErrNotFound
		}
		return nil, err
	}

	return mapPasteFromDB(update), nil
}

func (r *SqlcPasteRepository) Delete(ctx context.Context, userId int64, pasteUuid uuid.UUID) error {
	err := r.q.DeletePaste(ctx, dbgen.DeletePasteParams{
		Uuid:   pasteUuid,
		UserID: userId,
	})

	if err != nil {
		return err
	}

	return nil
}
