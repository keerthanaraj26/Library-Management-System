package handlers

import (
	"practice/database"
	"net/http"
	"practice/models"
	"encoding/json"
)
func GetStudentCatalogue(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, author, count, status FROM books")
	if err != nil {
		http.Error(w, "Error fetching books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Count, &b.Status)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"books": books})
}
