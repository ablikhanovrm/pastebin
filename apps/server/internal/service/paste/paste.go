package paste

import (
	"bytes"
	"context"
	"io"
	"time"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/ablikhanovrm/pastebin/internal/repository/cache"
	pasterepo "github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	cache     *cache.RedisCache
	log       zerolog.Logger
}

func NewPasteService(db *pgxpool.Pool, s3Storage *storage.Service, cache *cache.RedisCache, log zerolog.Logger) *Service {
	return &Service{db: db, s3Storage: s3Storage, cache: cache, log: log}
}

// repo helper
func (s *Service) repo(db dbgen.DBTX) *pasterepo.SqlcPasteRepository {
	return pasterepo.NewSqlcPasteRepository(db, s.log)
}

func (s *Service) Create(ctx context.Context, userId int64, in CreatePasteInput) (*paste.Paste, error) {
	repo := s.repo(s.db)

	newUuid := uuid.New()

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

	s3Key := newUuid.String()

	if err := s.s3Storage.Upload(ctx, s3Key, in.Content); err != nil {
		// cleanup db if s3 upload failed
		if delErr := repo.Delete(ctx, userId, newUuid); delErr != nil {
			log.Error().Err(delErr).Msg("failed cleanup paste after s3 fail")
		}
		return nil, ErrUploadFailed
	}

	go func() {
		if err := s.cache.SetPaste(ctx, opts); err != nil {
			log.Error().Err(err).Msg("failed save paste to cache after creating")
		}
	}()

	return createdPaste, nil
} // cache done

func (s *Service) GetByID(ctx context.Context, pasteUuid string, userId int64) (*paste.Paste, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return nil, err
	}

	res, err := s.cache.GetPaste(ctx, pasteUuid)

	if err != nil {
		s.log.Warn().Err(err).Msg("cache get failed")
	}

	if res != nil {
		return res, nil
	}

	res, err = repo.GetByID(ctx, userId, parsedUuid)

	if err != nil {
		return nil, err
	}

	go func() {
		if err = s.cache.SetPaste(context.Background(), res); err != nil {
			s.log.Warn().Err(err).Msg("failed to set paste in cache")
		}
	}()

	return res, nil
} // cache done

func (s *Service) GetContent(ctx context.Context, pasteUuid string, userId int64) (io.ReadCloser, int64, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return nil, 0, err
	}

	content, err := s.cache.GetPasteContent(ctx, pasteUuid)

	if err != nil {
		s.log.Warn().Err(err).Msg("cache get failed")
	}

	if len(content) > 0 {
		return io.NopCloser(bytes.NewReader(content)), int64(len(content)), nil
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
	buf := &bytes.Buffer{}
	tee := io.TeeReader(body, buf)

	// возвращаем tee как основной reader
	reader := io.NopCloser(tee)

	// кешируем ПОСЛЕ того как handler дочитает
	go func() {
		data, err := io.ReadAll(buf) // <--
		if err == nil && len(data) > 0 {
			_ = s.cache.SetPasteContent(context.Background(), pasteUuid, data)
		}
	}()

	if length == nil {
		return reader, 0, nil
	}
	return reader, *length, nil
} // cache done

func (s *Service) GetPastes(ctx context.Context, userId int64, cursor *time.Time, limit int32) ([]*paste.Paste, *time.Time, error) {
	repo := s.repo(s.db)

	return s.getFromCacheThenDB(
		ctx,
		cursor,
		limit,
		func() ([]*paste.Paste, error) {
			if cursor == nil {
				return repo.GetPastesFirstPage(ctx, pasterepo.GetPastesFirstPageParams{
					UserId: userId,
					Limit:  limit,
				})
			} else {
				return repo.GetPastesAfterCursor(ctx, pasterepo.GetPastesAfterCursorParams{
					UserId: userId,
					Cursor: *cursor,
					Limit:  limit,
				})
			}

		},
		func(ids []uuid.UUID) ([]*paste.Paste, error) {
			return repo.GetManyByIDs(ctx, userId, ids)
		},
	)
}

func (s *Service) GetMyPastes(ctx context.Context, userId int64, cursor *time.Time, limit int32) ([]*paste.Paste, *time.Time, error) {
	repo := s.repo(s.db)

	return s.getFromCacheThenDB(
		ctx,
		cursor,
		limit,
		func() ([]*paste.Paste, error) {
			if cursor == nil {
				return repo.GetPastesFirstPage(ctx, pasterepo.GetPastesFirstPageParams{
					UserId: userId,
					Limit:  limit,
				})
			} else {
				return repo.GetPastesAfterCursor(ctx, pasterepo.GetPastesAfterCursorParams{
					Cursor: *cursor,
					Limit:  limit,
					UserId: userId,
				})
			}

		},
		func(ids []uuid.UUID) ([]*paste.Paste, error) {
			return repo.GetManyByIDs(ctx, userId, ids)
		},
	)
}

func (s *Service) Delete(ctx context.Context, pasteUuid string, userId int64) error {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	err = repo.Delete(ctx, userId, parsedUuid)

	if err != nil {
		return err
	}

	err = s.s3Storage.Delete(ctx, pasteUuid)

	if err != nil {
		return err
	}

	go func() {
		_ = s.cache.DeletePaste(context.Background(), pasteUuid)
		_ = s.cache.InvalidatePasteLists(context.Background())
	}()

	return nil
}

func (s *Service) Update(ctx context.Context, pasteUuid string, userId int64, in UpdatePasteInput) error {
	parsedUuid, err := uuid.Parse(pasteUuid)
	if err != nil {
		return err
	}

	updated, err := s.repo(s.db).Update(ctx, userId, &paste.Paste{
		Uuid:       parsedUuid,
		Title:      in.Title,
		Syntax:     in.Syntax,
		Visibility: in.Visibility,
		MaxViews:   in.MaxViews,
		ExpiresAt:  in.ExpireAt,
	})

	if err != nil {
		return err
	}

	go func(p *paste.Paste) {
		err := s.cache.SetPaste(context.Background(), p)

		if err != nil {
			s.log.Warn().Err(err).Msg("failed to set paste in cache")
		}

		err = s.cache.InvalidatePasteLists(context.Background())

		if err != nil {
			s.log.Warn().Err(err).Msg("failed to invalidate paste lists")
		}
	}(updated)

	return nil
} // cache done

// helpers
func (s *Service) getFromCacheThenDB(
	ctx context.Context,
	cursor *time.Time,
	limit int32,
	fetchPage func() ([]*paste.Paste, error), // загрузка страницы из БД (если полный miss)
	fetchByIDs func([]uuid.UUID) ([]*paste.Paste, error), // загрузка конкретных paste по id
) ([]*paste.Paste, *time.Time, error) {

	var pastes []*paste.Paste

	// LIST CACHE
	pasteIDs, err := s.cache.GetPasteList(ctx, limit, cursor)
	if err != nil {
		s.log.Warn().Err(err).Msg("cache get list failed")
	}

	if len(pasteIDs) > 0 {
		founded, missIDs, err := s.cache.MgetPasteList(ctx, pasteIDs)
		if err != nil {
			s.log.Warn().Err(err).Msg("cache mget failed")
		}

		// MISS → DB
		var missPastes []*paste.Paste
		if len(missIDs) > 0 {
			parsed, err := parsePasteUuids(missIDs)
			if err == nil && len(parsed) > 0 {

				missPastes, err = fetchByIDs(parsed)
				if err != nil {
					s.log.Warn().Err(err).Msg("db load miss failed")
				}

				// прогрев кеша
				if len(missPastes) > 0 {
					go func(pastes []*paste.Paste) {
						if err := s.cache.MsetPasteList(ctx, pastes); err != nil {
							s.log.Warn().Err(err).Msg("cache set miss pastes failed")
						}
					}(missPastes)
				}
			}
		}

		// восстановление порядка
		resultMap := make(map[string]*paste.Paste, len(founded)+len(missPastes))

		for _, p := range founded {
			resultMap[p.Uuid.String()] = p
		}
		for _, p := range missPastes {
			resultMap[p.Uuid.String()] = p
		}

		for _, id := range pasteIDs {
			if p, ok := resultMap[id]; ok {
				pastes = append(pastes, p)
			}
		}

		if len(pastes) > 0 {
			next := pastes[len(pastes)-1].CreatedAt
			return pastes, &next, nil
		}

		s.log.Warn().Msg("cache list exists but empty result, fallback db")
	}

	// FULL DB FALLBACK
	pastes, err = fetchPage()
	if err != nil {
		return nil, nil, err
	}

	if len(pastes) == 0 {
		return pastes, nil, nil
	}

	// прогрев кеша
	go func(pastes []*paste.Paste, cur *time.Time, l int32) {
		if err := s.cache.MsetPasteList(ctx, pastes); err != nil {
			s.log.Warn().Err(err).Msg("cache set pastes failed")
		}

		ids := make([]string, 0, len(pastes))
		for _, p := range pastes {
			ids = append(ids, p.Uuid.String())
		}

		if err := s.cache.SetPasteList(ctx, ids, cur, l); err != nil {
			s.log.Warn().Err(err).Msg("cache set list failed")
		}

	}(pastes, cursor, limit)

	next := pastes[len(pastes)-1].CreatedAt
	return pastes, &next, nil
}
