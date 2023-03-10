package http_delivery

import (
	"net/http"
	"strconv"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/alfiankan/go-cqrs-blog/common"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
	httpResponse "github.com/alfiankan/go-cqrs-blog/transport/response"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type ArticleHTTPHandler struct {
	articleCommandUseCase domain.ArticleCommand
	articleQueryUseCase   domain.ArticleQuery
}

func NewArticleHTTPHandler(articleCommandUseCase domain.ArticleCommand, articleQueryUseCase domain.ArticleQuery) *ArticleHTTPHandler {
	return &ArticleHTTPHandler{articleCommandUseCase, articleQueryUseCase}
}

// CreateArticle godoc
// @Summary CreateArticle add/create articles
// @Description create new article
// @Tags articles
// @Accept json
// @Produce json
// @Param article body transport.CreateArticle true "Article detail"
// @Success 200
// @Router /articles [post]
func (handler *ArticleHTTPHandler) CreateArticle(c echo.Context) error {

	var reqBody transport.CreateArticle
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &httpResponse.HTTPBaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	if err := reqBody.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, &httpResponse.HTTPBaseResponse{
			Message: common.ValidationError.Error(),
			Data:    errVal,
		})
	}

	if err := handler.articleCommandUseCase.Create(c.Request().Context(), reqBody); err != nil {

		common.Log(common.LOG_LEVEL_ERROR, err.Error())

		return c.JSON(http.StatusInternalServerError, &httpResponse.HTTPBaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	return c.JSON(http.StatusCreated, &httpResponse.HTTPBaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    common.EmptyResponseData,
	})

}

// FindArticle godoc
// @Summary FindArticle find/get articles
// @Description Find articles, provided query param keyword and filter by author, data ordered by created time DESC
// @Tags articles
// @Accept json
// @Produce json
// @Param keyword query string false "search by keyword on title or body"
// @Param author query string false "filter by author"
// @Param page query string false "page result do tou large amount data, page from 1..n every page hold 50 articles, default page is 1"
// @Success 200
// @Router /articles [get]
func (handler *ArticleHTTPHandler) FindArticle(c echo.Context) error {

	keyword := c.QueryParam("keyword")
	author := c.QueryParam("author")
	pageParam := c.QueryParam("page")

	page, err := strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		page = 1
	}

	articles, err := handler.articleQueryUseCase.Get(c.Request().Context(), keyword, author, page)
	if err != nil {

		common.Log(common.LOG_LEVEL_ERROR, err.Error())

		return c.JSON(http.StatusInternalServerError, &httpResponse.HTTPBaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    common.EmptyResponseData,
		})
	}

	return c.JSON(http.StatusOK, &httpResponse.HTTPBaseResponse{
		Message: common.HttpSuccess,
		Data:    articles,
	})

}

func (handler *ArticleHTTPHandler) HandleRoute(e *echo.Echo) {
	e.POST("/articles", handler.CreateArticle)
	e.GET("/articles", handler.FindArticle)
}
