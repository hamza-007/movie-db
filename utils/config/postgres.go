package config

import (
	"os"
)

type postgreSQL struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Database string `validate:"required"`
	Password string `validate:"required"`
}

func (obj *postgreSQL) load() *postgreSQL {
	obj.Host = os.Getenv("POSTGRES_HOST")
	obj.Port = os.Getenv("POSTGRES_PORT")
	obj.User = os.Getenv("POSTGRES_USER")
	obj.Database = os.Getenv("POSTGRES_DB")
	obj.Password = os.Getenv("POSTGRES_PASSWORD")
	return obj
}
