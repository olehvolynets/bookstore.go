package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sethvargo/go-envconfig"
)

// Available params https://www.postgresql.org/docs/15/libpq-connect.html#LIBPQ-PARAMKEYWORDS
//
// Defaults:
//
//	user: OS user name
//	host: localhost
//	port: 5432 ?
//	dbname: same as user name
//	sslmode: prefered
type Config struct {
	Adapter string

	Scheme   string `env:"DB_SCHEME,default=postgresql"`
	User     string `env:"DB_USER,default=bookstore"`
	Password string `env:"DB_PASSWORD,default=bookstore_pass"`
	DbName   string `env:"DB_NAME,default=bookstore"`
	Host     string `env:"DB_HOST,default=127.0.0.1"`
	Port     uint   `env:"DB_PORT,default=5432"`

	Attrs url.Values
}

func (c Config) String() string {
	u := url.URL{
		Scheme:   c.Scheme,
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.DbName,
		RawQuery: c.Attrs.Encode(),
	}

	return u.String()
}

func NewConfig(ctx context.Context, adapter string) *Config {
	c := &Config{Adapter: adapter, Attrs: url.Values{}}

	envconfig.Process(ctx, c)

	return c
}
