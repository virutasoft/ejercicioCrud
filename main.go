package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	Id      int
	Name    string
	Suggest string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Client ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Client{}
	res := []Client{}
	for selDB.Next() {
		var id int
		var name, suggest string
		err = selDB.Scan(&id, &name, &suggest)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Suggest = suggest
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Client WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Client{}
	for selDB.Next() {
		var id int
		var name, suggest string
		err = selDB.Scan(&id, &name, &suggest)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Suggest = suggest
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Client WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Client{}
	for selDB.Next() {
		var id int
		var name, suggest string
		err = selDB.Scan(&id, &name, &suggest)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Suggest = suggest
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		suggest := r.FormValue("suggest")
		insForm, err := db.Prepare("INSERT INTO Client(name, suggest) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, suggest)
		log.Println("INSERT: Name: " + name + " | suggest: " + suggest)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		suggest := r.FormValue("suggest")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Client SET name=?, suggest=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, suggest, id)
		log.Println("UPDATE: Name: " + name + " | suggest: " + suggest)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Client WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8000")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	err := http.ListenAndServe(":8000", nil)
	log.Println(err.Error())
}
