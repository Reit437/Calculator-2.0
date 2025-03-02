package errors

import "errors"

var (
	ErrUnprocessableEntity = errors.New("Невалидные данные")
	ErrInternalServerError = errors.New("Что-то пошло не так")
	ErrNotFound            = errors.New("Нет такого выражения")
)
