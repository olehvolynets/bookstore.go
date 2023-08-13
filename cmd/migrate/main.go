package main

import (
	"log"

	"bookstore/internal/cli/migrate"
)

func main() {
	if err := migrate.Command.Execute(); err != nil {
		log.Fatal(err)
	}
}
