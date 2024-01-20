package synology

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Connection struct {
	client   *http.Client
	username string
	password string
	url      string
	sid      string
}

const PAGE_LIMIT = 100

func NewConnection(url string, username string, password string, insecureSkipVerify bool) Connection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}

	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	return Connection{
		username: username,
		password: password,
		url:      url,
		client:   &client,
	}
}

// https://github.com/zeichensatz/SynologyPhotosAPI

func (c *Connection) Login() error {
	body := fmt.Sprintf("api=SYNO.API.Auth&version=3&method=login&account=%s&passwd=%s", c.username, c.password)
	url := fmt.Sprintf("%s/webapi/auth.cgi", c.url)

	b := strings.NewReader(body)
	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}

	log.Printf("Performing login request to %s", url)
	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	parsedResp := &LoginResponse{}
	derr := json.NewDecoder(resp.Body).Decode(parsedResp)
	if derr != nil {
		return derr
	}

	if parsedResp.Success == true {
		c.sid = parsedResp.Data.Sid
		return nil
	}

	return errors.New("failed to extract sid from login to synology nas")
}

