package database

import (
	"bookstore/internal/books/model"
	"bookstore/internal/db"
)

type BooksInDB struct {
	db *db.DB
}

func (b *BooksInDB) FindById(id uint) (*model.Book, error) {
	return nil, nil
}
