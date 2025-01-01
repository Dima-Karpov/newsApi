package router

import (
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"newsApi/internal/domain"
	"strconv"
)

type getNewsResponse struct {
	Data        []domain.NewsList `json:"data"`
	CurrentPage int               `json:"currentPage"`
	Limit       int               `json:"limit"`
	TotalCount  int               `json:"totalCount"`
	TotalPages  int               `json:"totalPages"`
}

// @Summary Get news
// @Tags news
// @Description get news
// @Accept  json
// @Produce  json
//
//	@Param        page    query     number  false  "page"  Format(number)
//
// @Success 200 {object} getNewsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /news [get]
func (h *Handler) getNews(c *gin.Context) {
	// Получаем параметр page из query-параметров, eсли его нет, по умолчанию ставим 1
	page := c.DefaultQuery("page", "1")

	// Преобразуем парамтр page в число
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid page parameter")
		return
	}
	limit := 10

	news, totalCount, err := h.services.GetNews(pageInt, limit)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Вычисляем общее количество страниц
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	c.JSON(http.StatusOK, getNewsResponse{
		Data:        news,
		CurrentPage: pageInt,
		Limit:       limit,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
	})
}
