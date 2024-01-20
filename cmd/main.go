package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daspawnw/sync-google-photos/pkg/configuration"
	"github.com/daspawnw/sync-google-photos/pkg/googlephotos"
	"github.com/daspawnw/sync-google-photos/pkg/sync"
	"github.com/daspawnw/sync-google-photos/pkg/synology"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	c, err := configuration.LoadConfigurationFromFile(version, commit, date)
	if err != nil {
		log.Fatalf("Failed to load from configuration.yaml with error %v", err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	l := googlephotos.NewLogin(*c)
	client, err := l.PerformLogin(fmt.Sprintf("%s/oauth-login-state.json", path))
	if err != nil {
		log.Fatal(err)
	}

	if c.LoginOnly {
		log.Printf("Login to Google was successful and available at path %s/oauth-login-state.json", path)
		os.Exit(0)
	}

	s := synology.NewConnection(c.SynologyConf.Url, c.SynologyConf.Username, c.SynologyConf.Password, c.SynologyConf.InsecureHttpsConnection)
	lErr := s.Login()
	if lErr != nil {
		log.Fatalf("%v", err)
	}

	syncClient := sync.NewSync(&s, client, c.SynologyConf.UploadedTagID)
	syncErr := syncClient.Start()
	if syncErr != nil {
		log.Fatalf("%v", syncErr)
	}
}
