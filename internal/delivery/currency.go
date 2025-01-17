package delivery

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/crypto-service/internal/domain"
)

const (
	err500 = "500 Internal server error"
)

func (h *Handler) initAPI(router *gin.Engine) {
	currency := router.Group("/currency")

	currency.POST("/add", h.add)
}

func (h *Handler) add(c *gin.Context) {
	symbol := c.Query("symbol")

	if err := h.service.AddToTracking(symbol); err != nil {
		if errors.Is(err, domain.ErrSybmolIsEmpty) || errors.Is(err, domain.ErrSymbolDoesNotExists) {
			errResponse(c, http.StatusBadRequest, err.Error(), err.Error())

			return
		}

		errResponse(c, http.StatusInternalServerError, err.Error(), err500)

		return
	}

	newResponse(c, http.StatusOK, nil)
}
