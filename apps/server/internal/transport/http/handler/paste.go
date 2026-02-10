package handler

import (
	"errors"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetPastes(c *gin.Context) {

}

func (h *Handler) GetMyPastes(c *gin.Context) {

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

	if err != nil && errors.Is(paste.ErrNotFound, err) {
		c.JSON(404, gin.H{"error": "paste not found"})
		return
	}

	defer func() {
		_ = reader
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

}

func (h *Handler) DeletePaste(c *gin.Context) {

}

func (h *Handler) UpdatePaste(c *gin.Context) {

}
