package controller

import (
	"cms_sheets/lib"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type CallbackParams struct {
	Code  string `form:"code"`
	State string `form:"state"`
	Scope string `form:"scope"`
}

type OauthController struct{}

func (h OauthController) Callback(c *gin.Context) {
	var callbackPrams CallbackParams
	c.ShouldBindQuery(&callbackPrams)
	author := new(lib.GoogleApiAuth)
	author.SaveToken("token.json", callbackPrams.Code)

	location := url.URL{Path: "/test/show"}
	c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
}
