package app

import (
	"errors"
	"github.com/go-playground/validator/v10"
	myValidator "userSystem/pkg/validator"
)

func (g *Gin) ParseRequest(request interface{}) bool {
	if err := g.C.ShouldBind(request); err != nil {
		g.BadResponse(err.Error())
		return false
	}
	if err := myValidator.Validate.Struct(request); err != nil {
		var errStr string
		switch err.(type) {
		case validator.ValidationErrors:
			errStr = myValidator.Translate(err.(validator.ValidationErrors))
		default:
			errStr = errors.New("unknown error").Error()
		}
		g.BadResponse(errStr)
		return false
	}
	return true
}
