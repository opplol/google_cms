package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "embed"

	"golang.org/x/oauth2"
)

type auth_common_methods struct{}

func (au *auth_common_methods) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("erro:%v", err)
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err

}

// Retrieve a token, saves the token, then returns the generated client.
func (au *auth_common_methods) getClient(tokenfile string, config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := au.tokenFromFile(tokenfile)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}
