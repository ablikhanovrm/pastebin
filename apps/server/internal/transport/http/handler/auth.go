package handler

import (
	"net/http"

	"github.com/ablikhanovrm/pastebin/internal/service/auth"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	tokens, err := h.services.Auth.Login(ctx, auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		IP:        middleware.GetClientIP(c),
		UserAgent: middleware.GetUserAgent(c),
	})

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()}) // TODO error mapping
		return
	}

	secure := h.cfg.SecureCookies

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	resp := LoginResponse{
		AccessToken: tokens.AccessToken,
	}

	c.JSON(200, resp)
	return
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ip := middleware.GetClientIP(c)
	ua := middleware.GetUserAgent(c)

	tokens, err := h.services.Auth.Register(
		c.Request.Context(),
		auth.RegisterInput{
			Email:     req.Email,
			Password:  req.Password,
			IP:        ip,
			UserAgent: ua,
		},
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()}) // TODO error mapping
		return
	}

	secure := h.cfg.SecureCookies
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	resp := RegisterResponse{
		AccessToken: tokens.AccessToken,
	}

	c.JSON(200, resp)
}
