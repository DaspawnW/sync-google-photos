package googlephotos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daspawnw/sync-google-photos/pkg/configuration"
	"golang.org/x/oauth2"
)

type Login struct {
	oauthConf  *oauth2.Config
	appConf    configuration.Conf
	httpClient *http.Client
	ctx        context.Context
}

func NewLogin(c configuration.Conf) Login {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	oauthConf := &oauth2.Config{
		ClientID:     c.GooglePhotoConf.ClientID,
		ClientSecret: c.GooglePhotoConf.ClientSecret,
		Scopes:       []string{"openid", "profile", "https://www.googleapis.com/auth/photoslibrary.appendonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: c.GooglePhotoConf.LocalRedirectURL,
	}

	return Login{
		httpClient: httpClient,
		oauthConf:  oauthConf,
		appConf:    c,
		ctx:        ctx,
	}
}

func (l *Login) PerformLogin(filePath string) (*http.Client, error) {
	loginServerPort, err := l.appConf.GooglePhotoConf.LocalRedirectPort()
	if err != nil {
		return nil, fmt.Errorf("failed to load local redirect port from configuration with error %v", err)
	}

	loginServerPath, err := l.appConf.GooglePhotoConf.LocalRedirectPath()
	if err != nil {
		return nil, fmt.Errorf("failed to load local redirect path from configuration with error %v", err)
	}

	if l.fileAlreadyExists(filePath) {
		log.Printf("Found already existing login to google at path %s", filePath)
		tok, err := l.returnExistingFile(filePath)
		if err != nil {
			return nil, err
		}

		return l.oauthConf.Client(l.ctx, tok), nil
	}

	log.Printf("No existing login to google at path %s", filePath)
	url := l.oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("approval_prompt", "force"))

	dataChan := make(chan *string)
	LoginServer(dataChan, loginServerPort, *loginServerPath)

	log.Printf("Authentication URL: %s\n", url)

	code := <-dataChan

	tok, err := l.oauthConf.Exchange(l.ctx, *code)
	if err != nil {
		return nil, err
	}

	err = l.saveLogin(filePath, tok)
	if err != nil {
		return nil, err
	}

	return l.oauthConf.Client(l.ctx, tok), nil
}

func (l *Login) fileAlreadyExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (l *Login) returnExistingFile(filePath string) (*oauth2.Token, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var tok oauth2.Token
	uErr := json.Unmarshal(byteValue, &tok)
	if uErr != nil {
		return nil, uErr
	}

	return &tok, nil
}

func (l *Login) saveLogin(filePath string, tok *oauth2.Token) error {
	file, err := json.MarshalIndent(tok, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
