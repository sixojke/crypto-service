package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/crypto-service/pkg/logger"
)

type Response struct {
	Error    *errorResponse `json:"error,omitempty"`
	Response interface{}    `json:"response,omitempty"`
}

type errorResponse struct {
	Code int    `json:"code,omitempty"`
	Text string `json:"text,omitempty"`
}

func newResponse(c *gin.Context, statusCode int, response interface{}) {
	c.AbortWithStatusJSON(statusCode, Response{
		Response: response,
	})

	logger.Debugf("response: %v", response)
}

func errResponse(с *gin.Context, statusCode int, err, errResp string) {
	с.AbortWithStatusJSON(statusCode, Response{
		Error: &errorResponse{
			Code: statusCode,
			Text: errResp,
		},
	})

	logger.Warn(err)
}
