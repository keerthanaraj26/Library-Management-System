package models

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Count  int    `json:"count"`
	Status string `json:"Status"`
}
type Request struct {
	Student string `json:"student"`
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
