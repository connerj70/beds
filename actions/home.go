package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index2.html"))
}

func Index(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("index.html"))
}

func AboutHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("about.plush.html"))
}

func FAQHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("faq.plush.html"))
}
