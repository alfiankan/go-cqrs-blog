package http_delivery

import (
	"net/http"

	domain "github.com/alfiankan/go-cqrs-blog/article"
	"github.com/alfiankan/go-cqrs-blog/common"
	transport "github.com/alfiankan/go-cqrs-blog/transport/request"
	httpResponse "github.com/alfiankan/go-cqrs-blog/transport/response"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

type ArticleHTTPHandler struct {
	articleCommandUseCase domain.ArticleCommand
	articleQueryUseCase   domain.ArticleQuery
}

func NewArticleHTTPHandler(articleCommandUseCase domain.ArticleCommand, articleQueryUseCase domain.ArticleQuery) *ArticleHTTPHandler {
	return &ArticleHTTPHandler{articleCommandUseCase, articleQueryUseCase}
}

func (handler *ArticleHTTPHandler) CreateArticle(c echo.Context) error {

	var reqBody transport.CreateArticle
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &httpResponse.HTTPBaseResponse{
			Message: common.BadRequestError.Error(),
			Data:    nil,
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
		return c.JSON(http.StatusInternalServerError, &httpResponse.HTTPBaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
		})
	}

	return c.JSON(http.StatusCreated, &httpResponse.HTTPBaseResponse{
		Message: common.HttpSuccessCreated,
		Data:    nil,
	})

}

func (handler *ArticleHTTPHandler) FindArticle(c echo.Context) error {

	keyword := c.QueryParam("keyword")
	author := c.QueryParam("author")

	articles, err := handler.articleQueryUseCase.Get(c.Request().Context(), keyword, author)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &httpResponse.HTTPBaseResponse{
			Message: common.InternalServerError.Error(),
			Data:    nil,
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
