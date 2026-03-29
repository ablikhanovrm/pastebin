package handler

import (
	"fmt"
	"net/http"

	"github.com/ablikhanovrm/pastebin/internal/service/auth"
	"github.com/ablikhanovrm/pastebin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (h *Handler) Login(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("login failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Login req INPUT:", req.Email)
	fmt.Println("Login req INPUT:", req.Password)
	ctx := c.Request.Context()
	tokens, err := h.services.Auth.Login(ctx, auth.LoginInput{
		Email:     req.Email,
		Password:  req.Password,
		IP:        middleware.GetClientIP(c),
		UserAgent: middleware.GetUserAgent(c),
	})

	if err != nil {
		log.Warn().Err(err).Msg("login failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO error mapping
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

	c.JSON(http.StatusOK, resp)
	return
}

func (h *Handler) Register(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().
			Err(err).
			Str("email", req.Email).
			Msg("register failed")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "registration failed",
		})
		return
	}
	fmt.Println("Register req INPUT:", req.Email)
	fmt.Println("Register req INPUT:", req.Password)

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
		log.Warn().Err(err).Msg("register failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // TODO error mapping
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

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Logout(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	ctx := c.Request.Context()

	rt, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}

	err = h.services.Auth.Logout(ctx, rt)

	if err != nil {
		log.Warn().Err(err).Msg("logout failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	secure := h.cfg.SecureCookies
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}

func (h *Handler) RefreshToken(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	rt, err := c.Cookie("refresh_token")

	if err != nil || rt == "" {
		log.Warn().Err(err).Msg("refresh token failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}) // TODO error mapping
		return
	}

	tokens, err := h.services.Auth.Refresh(c.Request.Context(), rt, middleware.GetClientIP(c), middleware.GetUserAgent(c))

	if err != nil {
		log.Warn().Err(err).Msg("refresh token failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token failed"})
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

	c.JSON(http.StatusOK, resp)
}
