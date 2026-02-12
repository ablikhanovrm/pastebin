package paste

import (
	"context"
	"io"
	"time"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	pasterepo "github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PasteService interface {
	Create(ctx context.Context, u *paste.Paste) (*paste.Paste, error)
	GetByID(ctx context.Context, id uuid.UUID) (*paste.Paste, error)
	GetPastes(ctx context.Context, userId int64) ([]*paste.Paste, *time.Time, error)
	GetMyPastes(ctx context.Context, userId int64) ([]*paste.Paste, *time.Time, error)
	GetContent(ctx context.Context, pasteUuid string, userId int64) (io.ReadCloser, int64, error)
	Update(ctx context.Context, pasteUuid string, userId int64, in UpdatePasteInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	db        *pgxpool.Pool
	s3Storage *storage.Service
	log       zerolog.Logger
}

func NewPasteService(db *pgxpool.Pool, s3Storage *storage.Service, log zerolog.Logger) *Service {
	return &Service{db: db, s3Storage: s3Storage, log: log}
}

// repo helper
func (s *Service) repo(db dbgen.DBTX) *pasterepo.SqlcPasteRepository {
	return pasterepo.NewSqlcPasteRepository(db, s.log)
}

func (s *Service) Create(ctx context.Context, userId int64, in CreatePasteInput) (*paste.Paste, error) {
	repo := s.repo(s.db)

	newUuid := uuid.New()

	s3Key := newUuid.String()
	err := s.s3Storage.Upload(ctx, s3Key, in.Content)

	if err != nil {
		return nil, ErrUploadFailed
	}

	opts := &paste.Paste{
		Uuid:       newUuid,
		UserId:     userId,
		Title:      in.Title,
		Content:    nil,
		S3Key:      newUuid.String(),
		Syntax:     in.Syntax,
		Visibility: in.Visibility,
		MaxViews:   in.MaxViews,
		ExpiresAt:  in.ExpireAt,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	createdPaste, err := repo.Create(ctx, userId, opts)

	if err != nil {
		return nil, err
	}

	return createdPaste, nil
}

func (s *Service) GetByID(ctx context.Context, pasteUuid string, userId int64) (*paste.Paste, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return nil, err
	}

	res, err := repo.GetByID(ctx, userId, parsedUuid)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetContent(ctx context.Context, pasteUuid string, userId int64) (io.ReadCloser, int64, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return nil, 0, err
	}

	res, err := repo.GetByID(ctx, userId, parsedUuid)

	if err != nil {
		return nil, 0, err
	}

	if res == nil {
		return nil, 0, paste.ErrNotFound
	}

	body, length, err := s.s3Storage.Get(ctx, res.S3Key)

	if err != nil {
		return nil, 0, err
	}

	if length == nil {
		return body, 0, nil
	}

	return body, *length, nil
}

func (s *Service) GetPastes(ctx context.Context, userId int64, cursor *time.Time, limit int32) ([]*paste.Paste, *time.Time, error) {
	repo := s.repo(s.db)

	var pastes []*paste.Paste
	var err error

	if cursor == nil {
		pastes, err = repo.GetPastesFirstPage(ctx, pasterepo.GetPastesFirstPageParams{
			UserId: userId,
			Limit:  limit,
		})
	} else {
		pastes, err = repo.GetPastesAfterCursor(ctx, pasterepo.GetPastesAfterCursorParams{
			Cursor: *cursor,
			Limit:  limit,
			UserId: userId,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	if len(pastes) == 0 {
		return pastes, nil, nil
	}

	next := pastes[len(pastes)-1].CreatedAt

	return pastes, &next, nil

}

func (s *Service) GetMyPastes(ctx context.Context, userId int64, cursor *time.Time, limit int32) ([]*paste.Paste, *time.Time, error) {
	repo := s.repo(s.db)

	var pastes []*paste.Paste
	var err error

	if cursor == nil {
		pastes, err = repo.GetPastesFirstPage(ctx, pasterepo.GetPastesFirstPageParams{
			UserId: userId,
			Limit:  limit,
		})
	} else {
		pastes, err = repo.GetPastesAfterCursor(ctx, pasterepo.GetPastesAfterCursorParams{
			Cursor: *cursor,
			Limit:  limit,
			UserId: userId,
		})
	}

	if err != nil {
		return nil, nil, err
	}

	if len(pastes) == 0 {
		return pastes, nil, nil
	}

	next := pastes[len(pastes)-1].CreatedAt

	return pastes, &next, nil
}

func (s *Service) Delete(ctx context.Context, pasteUuid string, userId int64) error {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	err = repo.Delete(ctx, userId, parsedUuid)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, pasteUuid string, userId int64, in UpdatePasteInput) error {
	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return err
	}

	foundPaste, err := s.GetByID(ctx, pasteUuid, userId)

	if err != nil {
		return paste.ErrNotFound
	}

	if foundPaste == nil {
		return paste.ErrNotFound
	}

	err = s.repo(s.db).Update(ctx, userId, &paste.Paste{
		Uuid:       parsedUuid,
		Title:      in.Title,
		Syntax:     in.Syntax,
		Visibility: in.Visibility,
		MaxViews:   in.MaxViews,
		ExpiresAt:  in.ExpireAt,
	})

	if err != nil {
		return ErrUpdate
	}

	return nil
}
