package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	//Info logging - Important information
	Info *log.Logger
	//Warning logging - Concerning but not fatal error
	Warning *log.Logger
	//Error logging - Critical error
	Error *log.Logger
)
var mutex sync.Mutex //Mutex is used to define a critical section of code
var tpl *template.Template

func init() {
	//Creation of different logger type pointer values to support different logging levels
	file, err := os.OpenFile("logs/errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//Error when opening error.txt
		log.Fatalln("Failed to open error log file:", err)
	}
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//Path of templates
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	mutex.Lock() //Delete all User Sessions
	{
		db := OpenDB()
		defer db.Close()
		_, err := db.Exec(`DELETE FROM Sessions;`)
		if err != nil {
			Error.Println("Failed to delete SQL records:", err)
			log.Fatalln("Failed to delete in SQL records:", err)
		}
	}
	mutex.Unlock()

	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHTML)
	mux.HandleFunc("/signup", signupHTML)
	mux.HandleFunc("/login", loginHTML)
	mux.HandleFunc("/logout", logoutHTML)

	mux.HandleFunc("/explore", exploreHTML)
	mux.HandleFunc("/viewPost", viewPostHTML)
	mux.HandleFunc("/createPost", createPostHTML)
	mux.HandleFunc("/upload", uploadFileHTML)
	mux.HandleFunc("/contact", contactHTML)

	mux.HandleFunc("/weather", weatherHTML)

	mux.HandleFunc("/routes", routesHTML)
	mux.HandleFunc("/repeatCall", repeatedHome)
	mux.HandleFunc("/tracking", home)
	mux.HandleFunc("/location/{lat}/{long}", location)
	mux.HandleFunc("/viewRoute", viewRouteHTML)
	mux.HandleFunc("/createRoute", createRouteHTML)
	mux.HandleFunc("/heatmap", heatmapHTML)

	mux.HandleFunc("/account", accountHTML)
	mux.HandleFunc("/myPosts", myPostsHTML)
	mux.HandleFunc("/editPost", editPostHTML)
	mux.HandleFunc("/editUser", editUserHTML)

	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("templates/assets"))))

	err := http.ListenAndServeTLS(":5221", "D://GoLang//Projects//Go//src//GoLive//tls//cert.pem", "D://GoLang//Projects//Go//src//GoLive//tls//key.pem", mux)
	if err != nil {
		Error.Println("Unable to listen and serve TLS request: ", err)
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHTML(res http.ResponseWriter, req *http.Request) {
	userx := getUser(res, req)
	if userx.Username != "" {
		http.Redirect(res, req, "/explore", http.StatusSeeOther)
		Warning.Println("Already logged in.")
		return
	}
	tpl.ExecuteTemplate(res, "index.gohtml", nil)
}
