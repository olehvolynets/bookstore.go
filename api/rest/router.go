package rest

import (
	"github.com/uptrace/bunrouter"

	"bookstore/books"
)

func Routes(g *bunrouter.CompatGroup) {
	g.WithGroup("/books", func(bg *bunrouter.CompatGroup) {
		bg.GET("", books.HandleGetBooks)
	})
}
