package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:@/go_course?charset=utf8")

func main() {

	r := mux.NewRouter()
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{id}/view", ViewHandler)

	fmt.Println(http.ListenAndServe(":8080", r))
}

func ListPosts() []Post {
	rows, err := db.Query("SELECT * FROM posts")
	checkError(err)

	items := []Post{}
	for rows.Next() {
		post := Post{}

		rows.Scan(&post.Id, &post.Title, &post.Body)
		items = append(items, post)
	}
	return items
}

func GetPostById(id string) Post {
	row := db.QueryRow("SELECT * FROM posts WHERE id = ?", id)
	post := Post{}
	row.Scan(&post.Id, &post.Title, &post.Body)
	return post
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("templates/layout.html", "templates/list.html"))
	if err := template.ExecuteTemplate(w, "layout.html", ListPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	template := template.Must(template.ParseFiles("templates/layout.html", "templates/view.html"))
	if err := template.ExecuteTemplate(w, "layout.html", GetPostById(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
