package lib

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleApiAuthDrive struct {
	auth_common_methods
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
	tokFile := "drive_token.json"
	client, err := h.getClient(tokFile, config)
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
