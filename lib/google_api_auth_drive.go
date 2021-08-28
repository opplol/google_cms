package lib

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleApiAuthDrive struct {
	AuthUrl string `xorm: "Auth url string" form: "authurl" json: "authurl"`
	FileId  string `xorm: "cms file id string" form: "file id" json: "file id"`
}

//go:embed config/drive_credentials.json
var credentials_drive []byte

func (h GoogleApiAuthDrive) CredentialApi(docu_name string) (*GoogleApiAuthDrive, error) {
	ctx := context.Background()

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials_drive, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, err := h.getClient(config)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		h.AuthUrl = authURL
		return &h, err
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().Q(fmt.Sprintf("name = '%s'", docu_name)).
		Fields("files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			h.FileId = i.Id
		}
	}

	return &h, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func (h GoogleApiAuthDrive) getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "drive_token.json"
	tok, err := h.tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func (h GoogleApiAuthDrive) tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (h GoogleApiAuthDrive) SaveToken(path string, authCode string) {
	context.Background()

	fmt.Printf("%v\n", credentials_drive)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials_drive, drive.DriveMetadataReadonlyScope)
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
