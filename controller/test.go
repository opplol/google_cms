package controller

import (
	"cms_sheets/lib"
	"cms_sheets/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentController struct{}

func (h DocumentController) Show(c *gin.Context) {
	// test := new(model.DocuFirst)
	// test.GetData()

	documentModel := make([]model.DocuFirst, 0)

	valid_google_oauth_and_get_data(c, &documentModel)
	fmt.Printf("%#v\n", documentModel)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"data_model": documentModel,
	})
}

func valid_google_oauth_and_get_data(c *gin.Context, sheetModel *[]model.DocuFirst) {
	author := new(lib.GoogleApiAuth)
	authURL, err := author.Credential_api(sheetModel)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
