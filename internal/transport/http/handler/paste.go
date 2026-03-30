package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	pasteService "github.com/ablikhanovrm/pastebin/internal/service/paste"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (h *Handler) GetPastes(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	limit := int32(20)
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = int32(v)
		}
	}

	var cursor *time.Time
	if cur := c.Query("cursor"); cur != "" {
		t, err := time.Parse(time.RFC3339, cur)
		if err != nil {
			log.Warn().Err(err).Msg("parse cursor failed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor"})
			return
		}
		cursor = &t
	}

	items, nextCursor, err := h.services.Paste.GetPastes(c.Request.Context(), middleware.GetUserID(c), cursor, limit)
	if err != nil {
		log.Warn().Err(err).Msg("get pastes failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       items,
		"next_cursor": nextCursor,
	})
	return
}

func (h *Handler) GetMyPastes(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	limit := int32(20)
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = int32(v)
		}
	}

	var cursor *time.Time
	if cur := c.Query("cursor"); cur != "" {
		t, err := time.Parse(time.RFC3339, cur)
		if err != nil {
			log.Warn().Err(err).Msg("parse cursor failed")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor"})
			return
		}
		cursor = &t
	}

	items, nextCursor, err := h.services.Paste.GetPastes(c.Request.Context(), middleware.GetUserID(c), cursor, limit)
	if err != nil {
		log.Warn().Err(err).Msg("get pastes failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":       items,
		"next_cursor": nextCursor,
	})
	return
}

func (h *Handler) GetPaste(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingIDParam.Error()})
		return
	}

	userId := middleware.GetUserID(c)

	res, err := h.services.Paste.GetByID(c.Request.Context(), id, userId)

	if err != nil && errors.Is(paste.ErrNotFound, err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetPasteContent(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	id := c.Param("id")

	if id == "" {
		log.Warn().Msg("paste id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingIDParam.Error()})
		return
	}

	userId := middleware.GetUserID(c)

	reader, size, err := h.services.Paste.GetContent(c.Request.Context(), id, userId)

	if err != nil {
		if errors.Is(err, paste.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
			return
		}

		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	defer func() {
		_ = reader.Close()
	}()

	c.DataFromReader(
		200,
		size,
		"text/plain; charset=utf-8",
		reader,
		nil,
	)
	return
}

func (h *Handler) CreatePaste(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var opts CreatePasteRequest
	fmt.Println("PARAMS", c.Params)

	if err := c.ShouldBindJSON(&opts); err != nil {
		log.Warn().Err(err).Msg("invalid body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	input := pasteService.CreatePasteInput{
		Title:      opts.Title,
		Content:    opts.Content,
		Syntax:     opts.Syntax,
		Visibility: opts.Visibility,
		MaxViews:   opts.MaxViews,
		ExpireAt:   opts.ExpireAt,
	}

	res, err := h.services.Paste.Create(c.Request.Context(), middleware.GetUserID(c), input)

	if err != nil {
		log.Warn().Err(err).Msg("failed create paste")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create paste"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdatePaste(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var req UpdatePasteRequest

	pasteUuid := c.Param("id")
	if pasteUuid == "" {
		log.Warn().Msg("paste id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingIDParam.Error()})
	}

	if err := c.ShouldBindJSON(&req); err == nil {
		log.Warn().Err(err).Msg("invalid body")
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	updateOpts := pasteService.UpdatePasteInput{
		Title:      req.Title,
		Syntax:     req.Syntax,
		Visibility: req.Visibility,
		MaxViews:   req.MaxViews,
		ExpireAt:   req.ExpireAt,
	}

	err := h.services.Paste.Update(c.Request.Context(), pasteUuid, middleware.GetUserID(c), updateOpts)

	if err != nil {
		log.Warn().Err(err).Msg("failed update paste")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "paste updated"})
}

func (h *Handler) DeletePaste(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	pasteUuid := c.Param("id")

	err := h.services.Paste.Delete(c.Request.Context(), pasteUuid, middleware.GetUserID(c))

	if err != nil {
		log.Warn().Err(err).Msg("failed delete paste")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "paste deleted"})
}
