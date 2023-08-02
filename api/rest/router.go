package rest

import (
	"github.com/uptrace/bunrouter"

	"bookstore/internal/books"
)

func Routes(g *bunrouter.CompatGroup) {
	g.WithGroup("/books", func(bg *bunrouter.CompatGroup) {
		bg.GET("", books.HandleGetBooks)
	})
}
