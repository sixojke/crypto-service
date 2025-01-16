package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/crypto-service/internal/service"
)

type Handler struct {
	service *service.CurrencyService
}

func NewHandler(service *service.CurrencyService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init() *gin.Engine {
	// Create a new router
	router := gin.Default()

	router.GET("/ping", h.ping)

	h.initAPI(router)

	return router
}

// Test route to check server functionality
func (h *Handler) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
