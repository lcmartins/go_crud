package cli

import (
	"fmt"
	"gobooks/internal/service"
	"os"
	"strconv"
	"time"
)

type BookCLI struct {
	bookservice *service.BookService
}

func NewBookCli(s *service.BookService) *BookCLI {
	return &BookCLI{bookservice: s}
}

func (cli *BookCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: books <command> [arguments]")
		return
	}
	command := os.Args[1]

	switch command {
		case "search":
			if len(os.Args) < 3 {
				fmt.Println("Usage: books search <book title>")
			}
			bookName := os.Args[2]

			cli.SearchBooks(bookName)
		case "simulate":
			if len(os.Args) < 3 {
				fmt.Println("Similate: books simulate <book_id> <book_id> <book_id> ...")
			}
			bookids:= os.Args[2:]

			cli.SimulateReading(bookids)
	}

}

func(cli *BookCLI) SimulateReading(bookIdsString []string) {
	var booksIds []int 

	for _, idStr := range bookIdsString {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			fmt.Println("Invalid book id: ", idStr)
			continue
		}
		booksIds = append(booksIds, id)
	}

	responses := cli.bookservice.SimulateMultipleReadings(booksIds, 2*time.Second)

	for _, response:= range responses {
		fmt.Println(response)
	}
}

func(cli *BookCLI) SimulateReadingWithResponse(bookIdsString []string) []string {
	var booksIds []int 

	for _, idStr := range bookIdsString {
		id, err := strconv.Atoi(idStr)

		if err != nil {
			fmt.Println("Invalid book id: ", idStr)
			continue
		}
		booksIds = append(booksIds, id)
	}

	responses := cli.bookservice.SimulateMultipleReadings(booksIds, 2*time.Second)
	var result []string

	result = append(result, responses...)
	
	return result
}

func (cli *BookCLI) SearchBooks(name string) {
	books, err := cli.bookservice.SearchBooksByName(name)

	if err != nil {
		fmt.Println("Error searching books: ", err)
		return
	}

	if len(books) == 0 {
		fmt.Println("no books found")
		return
	}

	fmt.Printf("%d books found: \n", len(books))

	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Genre)
	}
}

func (cli *BookCLI) SearchBooksJSON(name string) ([]service.Book, error) {
	books, err := cli.bookservice.SearchBooksByName(name)

	if err != nil {
		fmt.Println("Error searching books: ", err)
		return nil, err
	}

	if len(books) == 0 {
		fmt.Println("no books found")
		return books, nil
	}

	fmt.Printf("%d books found: \n", len(books))

	return books, nil
}
