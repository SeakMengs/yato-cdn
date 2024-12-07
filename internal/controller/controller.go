package controller

import (
	appcontext "github.com/SeakMengs/yato-cdn/internal/app_context"
)

type baseController struct {
	app *appcontext.Application
}

type Controller struct {
	Index *IndexController
	File  *FileController
}

func newBaseController(app *appcontext.Application) *baseController {
	return &baseController{app: app}
}

func NewController(app *appcontext.Application) *Controller {
	bc := newBaseController(app)

	return &Controller{
		Index: &IndexController{baseController: bc},
		File:  &FileController{baseController: bc},
	}
}
