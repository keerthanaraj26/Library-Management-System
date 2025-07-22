package handlers

import (
	"encoding/json"
	"net/http"
	"practice/database"
	"practice/models"
	"html/template"
	"log"
)

func GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, student, title, status FROM requests WHERE status = 'pending'")
	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []models.Request
	for rows.Next() {
		var req models.Request
		err := rows.Scan(&req.Student, &req.ID, &req.Title, &req.Status)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"requests": requests})
}
func ApproveRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE requests SET status = 'approved' WHERE id = ?", data.ID)
	if err != nil {
		http.Error(w, "Failed to approve request", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Request approved successfully"})
}
func RequestsPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/approve_requests.html"))
	tmpl.Execute(w, nil)
}
func StudentCatalogue(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/student_catalogue.html"))
	tmpl.Execute(w, nil)
}

func ApplyForBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestData struct {
		BookID  int    `json:"book_id"`
		Student string `json:"student"`
	}

	var data RequestData

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("Decode error:", err)
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if data.Student == "" || data.BookID == 0 {
		http.Error(w, "Missing student or book ID", http.StatusBadRequest)
		return
	}

	var title, status string
	err := database.DB.QueryRow("SELECT title, status FROM books WHERE id = ?", data.BookID).Scan(&title, &status)
	if err != nil {
		log.Println("Book not found:", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if status != "active" {
		http.Error(w, "Book is not available", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO requests (student, id, title, status) VALUES (?, ?, ?, 'pending')",
		data.Student, data.BookID, title)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "Failed to apply", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Book request submitted!"})
}

func StudentHistory(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/student_history.html"))
	tmpl.Execute(w, nil)
}
func StudentHistoryHandler(w http.ResponseWriter, r *http.Request) {
	student := r.URL.Query().Get("email")
	if student == "" {
		http.Error(w, "Missing student email", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`SELECT id, title, status FROM requests WHERE student = ?`, student)
	if err != nil {
		http.Error(w, "Error fetching history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []models.Request

	for rows.Next() {
		var req models.Request
		req.Student = student
		if err := rows.Scan(&req.ID, &req.Title, &req.Status); err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			return
		}
		history = append(history, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
	})
}