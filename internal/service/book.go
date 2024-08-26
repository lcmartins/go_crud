package service

import (
	"database/sql"
	"fmt"
	"time"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Genre  string
}

func (book Book) GetFullBook() string {
	return book.Title + " by " + book.Author
}

type BookService struct {
	db *sql.DB
}

func NewBookService(db *sql.DB) *BookService {
	return &BookService{db: db}
}

func (service *BookService) CreateBook(book *Book) error {
	query := "Insert into books (title, author, genre) values(?,?,?)"
	result, err := service.db.Exec(query, book.Title, book.Author, book.Genre)

	if err != nil {
		return err
	}

	lastInsertedId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	book.ID = int(lastInsertedId)
	return nil
}

func (service *BookService) GetBooks() ([]Book, error) {
	query := "select id, title, author, genre from books"
	rows, err := service.db.Query(query)

	if err != nil {
		return nil, err
	}

	var books []Book

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}

func (service *BookService) GetBookById(id int) (*Book, error) {
	query := "select id, title, author, genre from books where id = ?"

	row := service.db.QueryRow(query, id)

	var book Book

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (service *BookService) UpdateBook(book Book) error {
	query := "update books set title=?, author=?, genre=? where id=?"

	_, err := service.db.Exec(query, &book.Title, &book.Author, &book.Genre, &book.ID)

	return err
}

func (service *BookService) DeleteBook(id int) error {
	query := "delete from books where id=?"

	_, err := service.db.Exec(query, id)

	return err
}

func (s *BookService) SimulateReading(bookId int, duration time.Duration, results chan<- string) {
	book, err := s.GetBookById(bookId)

	if err != nil || book == nil {
		results <- fmt.Sprintf("Book %d not found", bookId)
		return
	}

	time.Sleep(duration)
	results <- fmt.Sprintf("Book %s was read", book.Title)
}

func (s *BookService) SimulateMultipleReadings(bookIDs []int, duration time.Duration) []string {
	results := make(chan string, len(bookIDs))

	for _, id := range bookIDs {
		go func(bookId int) {
			s.SimulateReading(bookId, duration, results)
		}(id)
	}

	var responses []string
	for range bookIDs {
		responses = append(responses, <-results)
	}
	close(results)
	return responses
}

func (service *BookService) SearchBooksByName(name string) ([]Book, error) {
	query := "select id, title, author, genre from books where title like ?"

	rows, err := service.db.Query(query, "%"+name+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Genre)

		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}
