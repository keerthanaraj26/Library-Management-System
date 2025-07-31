package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	// "practice/database"
	"practice/handlers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)
type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestAdminLogin_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	// database.SetDB(db)

	email := "admin@example.com"
	password := "admin123"

	mock.ExpectQuery("SELECT password FROM users WHERE email = \\? AND role = \\?").
		WithArgs(email, "admin").
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(password))

	payload := LoginData{Email: email, Password: password}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login/admin", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.LoginHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestAdminLogin_InvalidPassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	// database.SetDB(db)

	email := "admin@example.com"
	password := "wrongpass"

	mock.ExpectQuery("SELECT password FROM users WHERE email = \\? AND role = \\?").
		WithArgs(email, "admin").
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("correctpass"))

	payload := LoginData{Email: email, Password: password}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login/admin", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.LoginHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAdminLogin_UserNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	// handlers.SetDB(db)

	email := "notfound@example.com"
	password := "any"

	mock.ExpectQuery("SELECT password FROM users WHERE email = \\? AND role = \\?").
		WithArgs(email, "admin").
		WillReturnError(sql.ErrNoRows)

	payload := LoginData{Email: email, Password: password}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/login/admin", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.LoginHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rec.Code)
	}
}

func TestAdminLogin_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/login/admin", nil)
	rec := httptest.NewRecorder()

	handlers.LoginHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestAdminLogin_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/login/admin", bytes.NewReader([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handlers.LoginHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}