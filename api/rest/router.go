package rest

import (
	"github.com/uptrace/bunrouter"

	"bookstore/books"
)

func Routes(g *bunrouter.Group) {
	g.WithGroup("/books", func(bg *bunrouter.Group) {
		bg.GET("", books.HandleGetBooks)
	})
}
