package main

import (
	"cms_sheets/controller"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func createHtmlRender() multitemplate.Renderer {
	renderer := multitemplate.NewRenderer()
	t := template.FuncMap{
		"safehtml": func(text string) template.HTML {
			return template.HTML(text)
		}, "evenOrOdd": func(index interface{}) (ret string) {
			if index.(int)%2 == 0 {
				ret = "even"
			} else {
				ret = "odd"
			}
			return
		},
	}
	baseLayout := []string{"view/base.tmpl", "view/index.tmpl", "view/common/base_head.tmpl",
		"view/common/side_bar.tmpl", "view/common/top_nav.tmpl",
		"view/common/include_js.tmpl", "view/common/footer.tmpl"}
	renderer.AddFromFilesFuncs("index", t, baseLayout...)
	renderer.AddFromFilesFuncs("content", t, "view/base.tmpl", "view/content.tmpl")
	renderer.AddFromFilesFuncs("event_list", t, append(baseLayout, "view/event/list.tmpl")...)

	return renderer
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// asset compile利用時のコード
	// views, err := loadTemplate()
	// if err != nil {
	// 	log.Fatalf("Error Load Template\n%#v", err)
	// }
	// r.SetHTMLTemplate(views)
	r.Use(static.Serve("/", static.LocalFile("./assets", false)))
	r.HTMLRender = createHtmlRender()

	r.GET("/panic", func(c *gin.Context) {
		panic("foo")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ohai")
	})

	cms_engine := r.Group("/pages")
	{
		document_controller := new(controller.DocumentController)
		cms_engine.GET("/show", document_controller.Show)
		cms_engine.GET("/:page", document_controller.Pages)
		cms_engine.GET("contents/:page", document_controller.PagesInContent)
	}
	oauth_engine := r.Group("/oauth")
	{
		oauth_controller := new(controller.OauthController)
		oauth_engine.GET("/callback", oauth_controller.Callback)
		oauth_engine.GET("/drive_callback", oauth_controller.DriveCallback)
	}
	test_engine := r.Group("/test")
	{
		test_controller := new(controller.TestController)
		test_engine.GET("/db_conn", test_controller.Show)
		test_engine.GET("/orm_test", test_controller.OrmTest)
	}
	event_engine := r.Group("/event")
	{
		event_controller := new(controller.EventController)
		event_engine.GET("/index", event_controller.Index)
		event_engine.GET("/push_me", event_controller.PushMe)
		event_engine.POST("/subscribe", event_controller.Subscribe)
	}
	r.Run(":8080")

}

// asset compile利用時のコード
// func loadTemplate() (*template.Template, error) {
// 	t := template.New("").Funcs(template.FuncMap{
// 		"safehtml": func(text string) template.HTML {
// 			return template.HTML(text)
// 		},
// 	})
// 	for name, file := range Assets.Files {
// 		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
// 			continue
// 		}
// 		h, err := ioutil.ReadAll(file)
// 		if err != nil {
// 			return nil, err
// 		}
// 		t, err = t.New(name).Parse(string(h))
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return t, nil
// }
