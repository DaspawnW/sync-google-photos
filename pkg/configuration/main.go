package configuration

import (
	"flag"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

func LoadConfigurationFromFile(version, commit, date string) (*Conf, error) {
	path := flag.String("file", "configuration.yaml", "Path to the configuration file")
	printVersion := flag.Bool("version", false, "Print version information")
	loginOnly := flag.Bool("login", false, "Only perform a login and generate the oauth-login-state.json file")

	flag.Parse()

	if *printVersion {
		log.Printf("Sync Google Photos application Version: %s on Commit: %s build at %s", version, commit, date)
		flag.PrintDefaults()
		os.Exit(0)
	}

	yamlFile, err := os.Open(*path)
	if err != nil {
		return nil, err
	}

	c := &Conf{}

	validate := validator.New()
	dec := yaml.NewDecoder(
		yamlFile,
		yaml.Validator(validate),
		yaml.Strict(),
	)

	decErr := dec.Decode(c)
	if decErr != nil {
		return nil, decErr
	}

	c.LoginOnly = *loginOnly

	return c, nil
}
