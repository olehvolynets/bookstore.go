package books

type Book struct{}

type BookService interface {
	GetAll() ([]Book, error)
	GetById() (Book, error)
}
