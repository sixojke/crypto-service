package delivery

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/sixojke/crypto-service/docs"
	"github.com/sixojke/crypto-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	config := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
	}

	router.Use(cors.New(config))

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", h.ping)

	h.initAPI(router)

	return router
}

// Test route to check server functionality
func (h *Handler) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
