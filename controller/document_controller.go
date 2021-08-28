package controller

import (
	"cms_sheets/lib"
	"cms_sheets/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentController struct{}

// GoogleSpred 連動テスト用アクション
func (h DocumentController) Show(c *gin.Context) {
	// test := new(model.DocuFirst)
	// test.GetData()

	documentModel := make([]model.DocuFirst, 0)

	valid_google_oauth_and_get_data(c, "test", &documentModel)
	fmt.Printf("%#v\n", documentModel)
	c.HTML(http.StatusOK, "/view/index.tmpl", gin.H{
		"data_model": documentModel,
	})
}

// GoogleDrive 連動テスト用アクション
func (h DocumentController) List(c *gin.Context) {
	// test := new(model.DocuFirst)
	// test.GetData()

	file_id := valid_google_oauth_and_get_list(c)
	fmt.Printf("%#v\n", file_id)
	c.JSON(http.StatusOK, gin.H{"file_id": file_id})
}

func (h DocumentController) Pages(c *gin.Context) {
	documentModel := make([]model.DocuFirst, 0)
	file_id := valid_google_oauth_and_get_list(c)
	fmt.Printf("%#v\n", file_id)

	valid_google_oauth_and_get_data(c, file_id, &documentModel)
	c.HTML(http.StatusOK, "/view/index.tmpl", gin.H{
		"data_model": documentModel,
	})
}

func valid_google_oauth_and_get_data(c *gin.Context, fileId string, sheetModel *[]model.DocuFirst) {
	author := new(lib.GoogleApiAuthSpread)
	authURL, err := author.CredentialApi(fileId, sheetModel)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, authURL)
	}
}
func valid_google_oauth_and_get_list(c *gin.Context) string {
	author := new(lib.GoogleApiAuthDrive)
	result, err := author.CredentialApi(c.Param("page"))
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, result.AuthUrl)
	}
	return result.FileId
}
