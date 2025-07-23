package main

import (
	"log"
	"net/http"
	"practice/database"
	"practice/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	database.Init()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.LoginPage)
	http.HandleFunc("/api/login/admin", handlers.AdminLoginHandler)
	http.HandleFunc("/api/login/student", handlers.StudentLoginHandler)
	http.HandleFunc("/admindashboard", database.AuthRequired("admin")(handlers.AdminDashboard))
	http.HandleFunc("/studentdashboard", database.AuthRequired("student")(handlers.StudentDashboard))
	http.HandleFunc("/register", database.AuthRequired("admin")(handlers.Register))
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/changepassword", database.AuthRequired("admin")(handlers.ChangePassword))
	http.HandleFunc("/api/change-password", handlers.ChangePasswordHandler)
	http.HandleFunc("/addbook", database.AuthRequired("admin")(database.AddPage))
	http.HandleFunc("/api/addbook", database.AddBookHandler)
	http.HandleFunc("/edit", database.AuthRequired("admin")(database.EditPage))
	http.HandleFunc("/api/edit", database.EditBookHandler)
	http.HandleFunc("/viewcatalogue", database.AuthRequired("admin")(database.ViewCatalogue))
	http.HandleFunc("/api/catalogue", handlers.GetCatalogue)
	http.HandleFunc("/archive", database.AuthRequired("admin")(database.ArchivePage))
	http.HandleFunc("/api/archive", database.ArchiveBookHandler)
	http.HandleFunc("/approverequests", database.AuthRequired("admin")(handlers.RequestsPage))
	http.HandleFunc("/api/requests", handlers.GetPendingRequests)
	http.HandleFunc("/api/approve", handlers.ApproveRequestHandler)
	http.HandleFunc("/studentcatalogue", database.AuthRequired("student")(handlers.StudentCatalogue))
	http.HandleFunc("/api/studentcatalogue", handlers.GetStudentCatalogue)
	http.HandleFunc("/api/apply", handlers.ApplyForBookHandler)
	http.HandleFunc("/studenthistory", database.AuthRequired("student")(handlers.StudentHistory))
	http.HandleFunc("/api/studenthistory", handlers.StudentHistoryHandler)

	log.Println("Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
