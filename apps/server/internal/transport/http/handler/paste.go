package handler

import "github.com/gin-gonic/gin"

func (h *Handler) GetPastes(ctx *gin.Context) {

}

func (h *Handler) GetMyPastes(ctx *gin.Context) {

}

func (h *Handler) GetPaste(ctx *gin.Context) {
	//body, size, err := h.services.Paste.Get(c.Request.Context(), key)
	//if err != nil {
	//	c.JSON(404, gin.H{"error": "not found"})
	//	return
	//}
	//defer body.Close()
	//
	//c.DataFromReader(
	//	200,
	//	size,
	//	"text/plain; charset=utf-8",
	//	body,
	//	nil,
	//)
}

func (h *Handler) CreatePaste(ctx *gin.Context) {

}

func (h *Handler) DeletePaste(ctx *gin.Context) {

}

func (h *Handler) UpdatePaste(ctx *gin.Context) {

}
