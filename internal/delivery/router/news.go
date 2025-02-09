package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// @Summary Get all news
// @Tags news
// @Description get news
// @Accept  json
// @Produce  json
//
//	@Param        page    	  query     number  false  "page"      Format(number)
//	@Param        fromDate    query     string  false  "fromDate"  Format(string)
//	@Param        toDate      query     string  false  "toDate"    Format(toDate)
//
// @Success 200 {object} getNewsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /news [get]
func (h *Handler) getAllNews(c *gin.Context) {
	// Получаем параметр page из query-параметров, eсли его нет, по умолчанию ставим 1
	page := c.DefaultQuery("page", "1")

	// Читаем параметры date из запроса
	fromDateStr := c.Query("fromDate")
	toDateStr := c.Query("toDate")

	var fromDate, toDate *string
	if fromDateStr != "" {
		fromDate = &fromDateStr
	}
	if toDateStr != "" {
		toDate = &toDateStr
	}

	// Преобразуем парамтр page в число
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid page parameter")
		return
	}
	limit := 10

	news, totalCount, err := h.services.GetNews(pageInt, limit, fromDate, toDate)
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

// @Summary Get news By ID
// @Tags news
// @Description get news by id
// @ID get-news-by-id
// @Accept  json
// @Produce  json
// @Param id path number true  "ID list"
// @Success 200 {object} domain.NewsList
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /news/{id} [get]
func (h *Handler) getNews(c *gin.Context) {
	// Получаем id из параметра маршрута
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id format")
		return
	}

	news, err := h.services.GetNew(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, news)
}
