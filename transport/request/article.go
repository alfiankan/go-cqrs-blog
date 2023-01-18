package transport

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateArticle struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (request CreateArticle) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Title, validation.Required, is.ASCII),
		validation.Field(&request.Author, validation.Required, is.Alphanumeric),
		validation.Field(&request.Body, validation.Required, is.ASCII),
	)
}
