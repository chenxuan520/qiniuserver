package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenInfo struct {
	Token string `json:"token"`
}

// Response 返回值
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	ErrHint string      `json:"err_hint,omitempty"`
}

// Success 成功
func Success(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
	})
}

// Error 错误
func Error(g *gin.Context, status int, data string) {
	g.JSON(status, Response{
		Code: status,
		Data: data,
	})
}

// Error 错误
// TODO
// record the wrong
func ErrorDetail(g *gin.Context, status int, data string, err error) {
	res := Response{
		Code: status,
		Data: data,
	}

	g.JSON(status, res)
}
