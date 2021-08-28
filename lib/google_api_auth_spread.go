package lib

import (
	"cms_sheets/model"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleApiAuthSpread struct {
	auth_common_methods
}

//go:embed config/spread_credentials.json
var credentials_spread []byte

func (h GoogleApiAuthSpread) CredentialApi(fileId string, sheetModel *[]model.DocuFirst) (string, error) {
	ctx := context.Background()

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials_spread, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	tokFile := "token.json"
	client, err := h.getClient(tokFile, config)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		return authURL, err
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	readRange := "sheet_cms1!A2:E"
	resp, err := srv.Spreadsheets.Values.Get(fileId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for index, row := range resp.Values {
			title := row[0].(string)
			body := row[1].(string)
			data := model.DocuFirst{Id: index, Title: title, Body: body}
			*sheetModel = append(*sheetModel, data)
		}
	}

	return "", nil

}

// Saves a token to a file path.
func (h GoogleApiAuthSpread) SaveToken(path string, authCode string) {
	context.Background()

	fmt.Printf("%#v\n", credentials_spread)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials_spread, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok, _ := config.Exchange(context.TODO(), authCode)
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
}
