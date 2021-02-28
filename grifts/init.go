package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/kolioDev/after_life/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
