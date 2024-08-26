package web

import (
	"encoding/json"
	"fmt"
	"gobooks/internal/cli"
	"gobooks/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type BookHanlders struct {
	service *service.BookService
	cli     *cli.BookCLI
}

func NewBookHanldlers(bookService *service.BookService, cli *cli.BookCLI) *BookHanlders {
	return &BookHanlders{service: bookService, cli: cli}
}

func (h *BookHanlders) GetBooks(w http.ResponseWriter, r *http.Request) {

	nameFromPath := r.URL.Query().Get("search")
	simulateFromPath := r.URL.Query().Get("simulate")

	if len(nameFromPath) > 0 {
		Search(nameFromPath, h, w)
		return
	}

	if len(simulateFromPath) > 0 {
		Simulate(h, simulateFromPath, w)
		return
	}

	GeAll(h, w)
}

func GeAll(h *BookHanlders, w http.ResponseWriter) {
	books, error := h.service.GetBooks()

	if error != nil {
		http.Error(w, "failed to get books", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func Simulate(h *BookHanlders, simulateFromPath string, w http.ResponseWriter) {
	booksRead := h.cli.SimulateReadingWithResponse(strings.Split(simulateFromPath, ","))
	fmt.Println(booksRead)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booksRead)
}

func Search(name string, h *BookHanlders, w http.ResponseWriter) {
	books, err := h.cli.SearchBooksJSON(name)

	if err != nil {
		http.Error(w, "failed to search", http.StatusInternalServerError)
	}

	if len(books) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHanlders) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		http.Error(w, "invalid payload request", http.StatusBadRequest)
		return
	}
	err = h.service.CreateBook(&book)
	fmt.Println("ERRI: ", err)

	if err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

}

func (h *BookHanlders) GetBookById(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromPath)

	if err != nil {
		http.Error(w, "failed to get the book id from request", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookById(id)

	if book == nil {
		http.Error(w, "failed to get the book", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "failed to get the book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHanlders) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromPath)

	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	var book service.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid payload request", http.StatusBadRequest)
		return
	}

	book.ID = id

	if err := h.service.UpdateBook(book); err != nil {
		http.Error(w, "failed to update the book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHanlders) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idFromPath := r.PathValue("id")

	id, err := strconv.Atoi(idFromPath)

	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBook(id); err != nil {
		http.Error(w, "failed to delete the book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (Handler *BookHanlders) SearchBooks(w http.ResponseWriter, r *http.Request) {
	nameFromPath := r.URL.Query().Get("search")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nameFromPath)
}
