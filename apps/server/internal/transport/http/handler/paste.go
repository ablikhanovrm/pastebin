package handler

import (
	"errors"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	pasteService "github.com/ablikhanovrm/pastebin/internal/service/paste"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPastes(c *gin.Context) {
	pastes, err := h.services.Paste.GetPastes(c.Request.Context(), middleware.GetUserID(c))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, pastes)
	return
}

func (h *Handler) GetMyPastes(c *gin.Context) {
	pastes, err := h.services.Paste.GetMyPastes(c.Request.Context(), middleware.GetUserID(c))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, pastes)
	return
}

func (h *Handler) GetPaste(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": ErrMissingIDParam.Error()})
		return
	}

	userId := middleware.GetUserID(c)

	res, err := h.services.Paste.GetByID(c.Request.Context(), id, userId)

	if err != nil && errors.Is(paste.ErrNotFound, err) {
		c.JSON(404, gin.H{"error": "paste not found"})
		return
	}

	c.JSON(200, res)
}

func (h *Handler) GetPasteContent(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(400, gin.H{"error": ErrMissingIDParam.Error()})
		return
	}

	userId := middleware.GetUserID(c)

	reader, size, err := h.services.Paste.GetContent(c.Request.Context(), id, userId)

	if err != nil {
		if errors.Is(err, paste.ErrNotFound) {
			c.JSON(404, gin.H{"error": "paste not found"})
			return
		}

		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	defer func() {
		_ = reader.Close()
	}()

	// TODO: add return data from cache(redis)
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
	var opts CreatePasteRequest

	if err := c.ShouldBindJSON(&opts); err == nil {
		c.JSON(400, gin.H{"error": "invalid body"})
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
		c.JSON(500, gin.H{"error": "failed create paste"})
		return
	}

	c.JSON(200, res)
}

func (h *Handler) UpdatePaste(c *gin.Context) {
	var req UpdatePasteRequest

	pasteUuid := c.Param("id")
	if pasteUuid == "" {
		c.JSON(400, gin.H{"error": ErrMissingIDParam.Error()})
	}

	if err := c.ShouldBindJSON(&req); err == nil {
		c.JSON(400, gin.H{"error": ErrInvalidJSON})
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "paste updated"})
}

func (h *Handler) DeletePaste(c *gin.Context) {
	pasteUuid := c.Param("id")

	err := h.services.Paste.Delete(c.Request.Context(), pasteUuid, middleware.GetUserID(c))

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"message": "paste deleted"})
}
