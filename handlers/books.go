package handlers

import (
	"encoding/json"
	"net/http"
	"practice/database"
	"practice/models"
)

func GetCatalogue(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, author, count FROM books")
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Count); err != nil {
			http.Error(w, "Error reading row", http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
