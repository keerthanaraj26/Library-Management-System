package handlers

import (
	"encoding/json"
	"net/http"
	"practice/database"
	"practice/models"
	"html/template"
	"github.com/go-chi/chi/v5"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, nil)
}
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/admin_dashboard.html"))
	tmpl.Execute(w, nil)
}
func StudentDashboard(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/student_dashboard.html"))
	tmpl.Execute(w, nil)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	role := chi.URLParam(r, "role")
	var creds models.Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var user models.User
	var query string

	if role == "admin" {
		query = "SELECT email, password FROM users WHERE email = ? AND password = ?"
	} else if role == "student" {
		query = "SELECT email, password FROM students WHERE email = ? AND password = ?"
	} else {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	err = database.DB.QueryRow(query, creds.Email, creds.Password).Scan(&user.Email, &user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	login(w, r, "admin")
}

func StudentLoginHandler(w http.ResponseWriter, r *http.Request) {
	login(w, r, "student")
}

func login(w http.ResponseWriter, r *http.Request, role string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data LoginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var storedPass string
	err := database.DB.QueryRow("SELECT password FROM users WHERE email = ? AND role = ?", data.Email, role).Scan(&storedPass)
	if err != nil || storedPass != data.Password {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: data.Email,
		Path:  "/",
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
	})
}
func Register(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	_, err := database.DB.Exec("INSERT INTO users (email, password, role) VALUES (?, ?, 'student')", data.Email, data.Password)
	if err != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/change_password.html"))
	tmpl.Execute(w, nil)
}
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		Role        string `json:"role"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	var storedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE email=? AND role=?", data.Email, data.Role).Scan(&storedPassword)
	if err != nil || storedPassword != data.OldPassword {
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		return
	}

	_, err = database.DB.Exec("UPDATE users SET password=? WHERE email=? AND role=?", data.NewPassword, data.Email, data.Role)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Path:     "/",
		MaxAge:   -1, 
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}