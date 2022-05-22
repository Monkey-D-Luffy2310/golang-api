package service

import (
	"fmt"
	"golang_api/dto"
	"golang_api/entity"
	"golang_api/repository"
	"log"

	"github.com/mashingan/smapping"
)

type BookService interface {
	Insert(book dto.BookCreateDTO) entity.Book
	Update(book dto.BookUpdateDTO) entity.Book
	Delete(book entity.Book)
	FindAll() []entity.Book
	FindByID(bookID uint64) entity.Book
	IsAllowedEdit(userID string, bookID uint64) bool
}

type bookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepo,
	}
}

func (b *bookService) Insert(book dto.BookCreateDTO) entity.Book {
	bookCreate := entity.Book{}
	err := smapping.FillStruct(&bookCreate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	createdBook := b.bookRepository.InsertBook(bookCreate)
	return createdBook
}

func (b *bookService) Update(book dto.BookUpdateDTO) entity.Book {
	bookUpdate := entity.Book{}
	err := smapping.FillStruct(&bookUpdate, smapping.MapFields(&book))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}
	updatedBook := b.bookRepository.UpdateBook(bookUpdate)
	return updatedBook
}

func (b *bookService) Delete(book entity.Book) {
	b.bookRepository.DeleteBook(book)
}

func (b *bookService) FindAll() []entity.Book {
	return b.bookRepository.FindAllBooks()
}

func (b *bookService) FindByID(bookID uint64) entity.Book {
	return b.bookRepository.FindBookByID(bookID)
}

func (b *bookService) IsAllowedEdit(userID string, bookID uint64) bool {
	book := b.bookRepository.FindBookByID(bookID)
	id := fmt.Sprintf("%v", book.UserID)
	return id == userID
}
