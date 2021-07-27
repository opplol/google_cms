package main

import (
	"cms_sheets/controller"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("foo")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ohai")
	})

	r.SetFuncMap(template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	})

	r.LoadHTMLGlob("view/*")
	test_engine := r.Group("/test")
	{
		document_controller := new(controller.DocumentController)
		test_engine.GET("/show", document_controller.Show)
	}
	oauth_engine := r.Group("/oauth")
	{
		oauth_controller := new(controller.OauthController)
		oauth_engine.GET("/callback", oauth_controller.Callback)
	}
	r.Run(":8080")

}
