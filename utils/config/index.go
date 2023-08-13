package config

import (
	"errors"
	"io/fs"
	"log"
	"strings"
	"sync"

	viper "github.com/spf13/viper"
)

/*============================================================================*/
/*=====*                            Accessor                            *=====*/
/*============================================================================*/

func Environment() string { return viper.GetString("ENVIRONMENT") }

func IsProduction() bool { return Environment() == "production" }

func IsStaging() bool { return Environment() == "staging" }

func IsDemo() bool { return Environment() == "demo" }

func IsDevelopment() bool { return Environment() == "development" }

func IsOnline() bool { return IsProduction() || IsStaging() || IsDemo() }

func IsTest() bool { return Environment() == "test" }

func Loki() loki {
	setup().lokiOnce.Do(func() { setup().lokiConfig.load() })
	return setup().lokiConfig
}

func PostgreSQL() postgreSQL {
	setup().pgOnce.Do(func() { setup().pgConfig.load() })
	return setup().pgConfig
}

/*============================================================================*/
/*=====*                            Container                           *=====*/
/*============================================================================*/

var once sync.Once
var config *container

type container struct {
	// Loki
	lokiOnce   sync.Once
	lokiConfig loki

	// PostgreSQL
	pgOnce   sync.Once
	pgConfig postgreSQL
}

func setup() *container {
	once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

		// Load from ENV
		viper.AutomaticEnv()
		err := viper.ReadInConfig()
		if err != nil {
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				// Do nothing
			} else if errors.Is(err, fs.ErrNotExist) {
				// Do nothing
			} else if err != nil {
				log.Fatalf("Fatal error loading ENV: %v", err)
			}
		}

		// Load from .env
		viper.SetConfigType("env")
		viper.SetConfigFile(".env")
		err = viper.MergeInConfig()
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			// Do nothing
		} else if errors.Is(err, fs.ErrNotExist) {
			// Do nothing
		} else if err != nil {
			log.Fatalf("Fatal error loading .env: %v", err)
		}

		// Populate config
		config = &container{}
	})
	return config
}
