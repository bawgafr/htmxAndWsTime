package boiler

import (
	"io"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const cookieName = "sessionId"

type Templates struct {
	template *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func newTemplate(embededStatic fs.FS) *Templates {
	return &Templates{
		//		template: template.Must(template.ParseGlob("views/*.html")),
		template: template.Must(template.ParseFS(embededStatic, "views/*.html")),
	}
}

func NewEcho(embededStatic fs.FS) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate(embededStatic)
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "static",
		Filesystem: http.FS(embededStatic),
	}))
	e.Use(middlewareReadSessionId)
	e.Static("/static/images", "images")
	e.Static("/static/css", "css")
	return e
}

func middlewareReadSessionId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var id string
		cookie, ok := getCookie(c)
		if !ok {
			id = uuid.NewString()
			setCookie(c, id)
		} else {
			id = cookie.Value
		}

		c.Set("sessId", id)

		return next(c)
	}
}
