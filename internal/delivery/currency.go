package delivery

import (
	"errors"
	"net/http"
	"strconv"

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
	currency.DELETE("/remove", h.removeCurrency)
	currency.GET("/price", h.getPrice)
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

// @Summary Remove currency from tracking
// @Tags currency
// @Description Removes a currency from the tracking list.
// @ModuleID removeCurrency
// @Accept json
// @Produce json
// @Param symbol query string true "Currency symbol (e.g., BTCUSDT)"
// @Success 200 {object} nil "Currency added successfully"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /currency/remove [delete]
func (h *Handler) removeCurrency(c *gin.Context) {
	symbol := c.Query("symbol")

	if err := h.service.RemoveFromTracking(symbol); err != nil {
		errResponse(c, http.StatusInternalServerError, err.Error(), err500)

		return
	}

	newResponse(c, http.StatusOK, nil)
}

// @Summary Get price of a currency at a specific timestamp
// @Tags currency
// @Description Retrieves the price of a specific currency at a given timestamp.
// @ModuleID getPrice
// @Accept json
// @Produce json
// @Param symbol query string true "Currency symbol (e.g., BTCUSDT)"
// @Param timestamp query string true "Currency symbol (e.g., BTCUSDT)"
// @Success 200 {object} Response "Price retrieved successfully"
// @Failure 400 {object} Response "Response "Bad Request (e.g., invalid symbol)"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /currency/price [get]
func (h *Handler) getPrice(c *gin.Context) {
	symbol := c.Query("symbol")
	timestamp := c.Query("timestamp")

	timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		errResponse(c, http.StatusBadRequest, domain.ErrInvalidTimestamp.Error(), domain.ErrInvalidTimestamp.Error())

		return
	}

	if symbol == "" {
		errResponse(c, http.StatusBadRequest, domain.ErrSybmolIsEmpty.Error(), domain.ErrSybmolIsEmpty.Error())

		return
	}

	price, err := h.service.GetPriceByTimestamp(symbol, timestampInt64)
	if err != nil {
		if errors.Is(err, domain.ErrNoDataOnThisCurrency) {
			errResponse(c, http.StatusBadRequest, err.Error(), err.Error())

			return
		}

		errResponse(c, http.StatusInternalServerError, err.Error(), err500)

		return
	}

	newResponse(c, http.StatusOK, Response{
		Response: price,
	})
}
