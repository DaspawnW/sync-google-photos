package configuration

import (
	"net/url"
	"strconv"
)

type Conf struct {
	SynologyConf    SynologyConf    `yaml:"synologyConfiguration" validate:"required"`
	GooglePhotoConf GooglePhotoConf `yaml:"googlePhotoConfiguration" validate:"required"`
	LoginOnly       bool
}

type SynologyConf struct {
	Url                     string `yaml:"url" validate:"required"`
	Username                string `yaml:"username" validate:"required"`
	Password                string `yaml:"password" validate:"required"`
	InsecureHttpsConnection bool   `yaml:"insecureHttpsConnection"`
	UploadedTagID           int    `yaml:"uploadedTagID" validate:"required"`
}

type GooglePhotoConf struct {
	ClientID         string `yaml:"clientID" validate:"required"`
	ClientSecret     string `yaml:"clientSecret" validate:"required"`
	LocalRedirectURL string `yaml:"localRedirectURL" validate:"required"`
}

func (g GooglePhotoConf) LocalRedirectPort() (int, error) {
	u, err := url.Parse(g.LocalRedirectURL)
	if err != nil {
		return 0, err
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return 0, err
	}

	return port, nil
}

func (g GooglePhotoConf) LocalRedirectPath() (*string, error) {
	u, err := url.Parse(g.LocalRedirectURL)
	if err != nil {
		return nil, err
	}

	return &u.Path, nil
}
