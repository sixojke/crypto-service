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
	api := router.Group("/api")

	currency := api.Group("/currency")
	currency.POST("/add", h.addCurrency)
}

// @Summary Add currency to tracking
// @Tags currency
// @Description Adds a currency to the tracking list.
// @ModuleID addCurrency
// @Accept json
// @Produce json
// @Param symbol query string true "Currency symbol (e.g., BTCUSDT)"
// @Success 200 {object} nil "Currency added successfully"
// @Failure 400 {object} Response "Bad Request (e.g., invalid symbol)"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /currency/add [post]
func (h *Handler) addCurrency(c *gin.Context) {
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
