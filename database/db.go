package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"practice/models"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Count  int    `json:"count"`
}

var DB *sql.DB

func Init() {
	var err error
	dsn := "root:test@123@tcp(127.0.0.1:3306)/demo1"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Unable to reach database:", err)
	}

}
func AddPage(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/add_book.html"))
	tmpl.Execute(w, nil)
}

func EditPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/edit_book.html"))
	tmpl.Execute(w, nil)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// http.Error(w, "Missing book ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book *models.Book
	err = DB.QueryRow("SELECT id, title, author, count FROM books WHERE id = ?", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Count)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

}

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data Book
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Received Data: %+v\n", data)

		_, err = DB.Exec(
			`INSERT INTO books (id, title, author, count) VALUES (?, ?, ?, ?)`,
			data.ID, data.Title, data.Author, data.Count,
		)
		if err != nil {
			http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "Data received successfully!"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func EditBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var editdata models.Book
	err := json.NewDecoder(r.Body).Decode(&editdata)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE books SET title=?, author=?, count=? WHERE id=?`
	result, err := DB.Exec(query, editdata.Title, editdata.Author, editdata.Count, editdata.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Book updated successfully!",
	})
}

func ViewCatalogue(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/view_catalogue.html")
	t.Execute(w, nil)
}
func ArchivePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/archive_book.html"))
	tmpl.Execute(w, nil)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// http.Error(w, "Missing book ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	err = DB.QueryRow("SELECT id, title, author, count FROM books WHERE id = ?", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Count)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

}
func ArchiveBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID int `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE books SET status = 'archived' WHERE id = ?`
	result, err := DB.Exec(query, data.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Book archived successfully!",
	})
}

