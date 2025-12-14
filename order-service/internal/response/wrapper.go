package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
    StatusCode int         `json:"statusCode"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data"`
}

type AppHandler func(c *gin.Context) (int, string, interface{})

func Wrap(h AppHandler) gin.HandlerFunc {
    return func(c *gin.Context) {
        code, msg, data := h(c)
        c.JSON(code, APIResponse{
            StatusCode: code,
            Message:    msg,
            Data:       data,
        })
    }
}
