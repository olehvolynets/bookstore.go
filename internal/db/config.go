package db

import (
	"fmt"
	"net/url"
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
	Scheme   string `env:"DB_SCHEME, default=postgresql"`
	User     string `env:"DB_USER, default=bookstore"`
	Password string `env:"DB_PASSWORD, default=bookstore_pass"`
	DbName   string `env:"DB_NAME, default=bookstore"`
	Host     string `env:"DB_HOST, default=127.0.0.1"`
	Port     uint   `env:"DB_PORT, default=5432"`
	SSLMode  string `env:"DB_SSLMODE, default=disable"`
}

func (c Config) String() string {
	params := url.Values{}
	params.Add("sslmode", c.SSLMode)

	u := url.URL{
		Scheme:   c.Scheme,
		User:     url.UserPassword(c.User, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.DbName,
		RawQuery: params.Encode(),
	}

	return u.String()
}
