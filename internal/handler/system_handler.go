package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SystemHandler struct {
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
	}
}


type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

type PingErrorResponse struct {
	Error string `json:"error" example:"Internal Server Error"`
}

// PingHandler godoc
// @Summary      Ping応答を確認
// @Description  Ping応答を確認するためのハンドラーです。そもそもswaggerが動いていればpongは絶対返ってくるはず。
// @Tags         system
// @Produce      json
// @Success      200  {object}  PingResponse
// @Failure      500  {object}  PingErrorResponse
// @Router       /ping [get]
func (h *SystemHandler) PingHandler(c *gin.Context) {
	// 疎通確認用のハンドラー
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}