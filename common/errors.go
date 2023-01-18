package common

import "errors"

var (
	ValidationError     = errors.New("Request format not valid")
	InternalServerError = errors.New("Something wrong in server")
	BadRequestError     = errors.New("Bad malformed request")
)
