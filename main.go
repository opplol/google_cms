package main

import (
	"cms_sheets/controller"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

	views, err := loadTemplate()
	if err != nil {
		log.Fatalf("Error Load Template\n%#v", err)
	}
	r.SetHTMLTemplate(views)

	test_engine := r.Group("/pages")
	{
		document_controller := new(controller.DocumentController)
		test_engine.GET("/show", document_controller.Show)
		// test_engine.GET("/list", document_controller.List)
		test_engine.GET("/:page", document_controller.Pages)
	}
	oauth_engine := r.Group("/oauth")
	{
		oauth_controller := new(controller.OauthController)
		oauth_engine.GET("/callback", oauth_controller.Callback)
		oauth_engine.GET("/drive_callback", oauth_controller.DriveCallback)
	}
	r.Run(":8080")

}

func loadTemplate() (*template.Template, error) {
	t := template.New("").Funcs(template.FuncMap{
		"safehtml": func(text string) template.HTML {
			return template.HTML(text)
		},
	})
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
