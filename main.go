package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Data struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var datas []Data

const fileName = "data.json"

func main() {
	http.HandleFunc("/", Handler)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/registration", Registration)
	http.HandleFunc("/loginPage", LoginPage)

	log.Println("server working at: http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "ERROR WHILE PARSING THE FILE: ", http.StatusBadRequest)
		return
	}
	temp.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	data := Data{
		Name:     r.FormValue("name"),
		ID:       id,
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	datas = append(datas, data)
	Saver()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("template/registration.html")
	if err != nil {
		http.Error(w, "ERROR WHILE PARSING THE FILE: ", http.StatusBadRequest)
	}
	temp.Execute(w, nil)
}

func Saver() {
	data, err := json.MarshalIndent(datas, "", "")
	if err != nil {
		log.Fatalln("ERROR WHILE MARSHALLING: ", err)
		return
	}
	err = os.WriteFile(fileName, data, 0o644)
	if err != nil {
		log.Fatalln("ERROR WHILE WRITING THE FILE", err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	name := r.FormValue("name")
	pass := r.FormValue("password")

	for _, data := range datas {
		if data.Name == name && data.Password == pass {
			fmt.Fprintf(w, "welcome: %s", name)
		}
	}
	http.Error(w, "INVALID USER NAME OR PASSWORD: ", http.StatusUnauthorized)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("template/login.html")
	if err != nil {
		http.Error(w, "ERROR WHILE PARSING THE FILE: ", http.StatusBadRequest)
	}
	temp.Execute(w, nil)
}
