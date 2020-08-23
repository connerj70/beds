package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// AuthIndex default implementation.
func AuthIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/index.html"))
}

func AuthSignIn(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/signin.html"))
}

func AuthRegister(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/register.html"))
}
