package config

import (
	"fmt"
	"log"
	"net/url"

	viper "github.com/spf13/viper"
)

type loki struct {
	url string `validate:"omitempty"`
	jwt string `validate:"omitempty"`
}

func (loki) namespace() string         { return "Loki" }
func (obj loki) key(key string) string { return fmt.Sprintf("%s_%s", obj.namespace(), key) }

func (obj *loki) load() {
	obj.url = viper.GetString(obj.key("URL"))
	obj.jwt = viper.GetString(obj.key("JWT"))
}

func (obj loki) Enabled() bool {
	if obj.url == "" {
		return false
	}
	_, err := url.Parse(obj.url)
	return err == nil
}

func (obj loki) URL() string {
	// Parse server URL
	u, err := url.Parse(obj.url)
	if err != nil {
		log.Fatal(err)
	}

	if obj.jwt != "" {
		return fmt.Sprintf("%s://api_key:%s@%s", u.Scheme, obj.jwt, u.Host)
	}
	return obj.url
}
