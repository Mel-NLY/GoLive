package main

import (
	"database/sql"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt" //go get golang.org/x/crypto/bcrypt
)

//User is exported
type User struct {
	Username string
	Password []byte
	First    string
	Last     string
	Email    string
}

func loginHTML(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		//User is already signed in
		http.Redirect(res, req, "/explore", http.StatusSeeOther)
		Warning.Println("Already logged in. Unable to log in again.")
		return
	}

	//Process form submission
	if req.Method == http.MethodPost { //Check if the form has been submitted
		username := req.FormValue("username_input")
		if x, _ := regexp.MatchString("^[a-zA-Z0-9]+$", req.FormValue("username_input")); !x || username == "" {
			http.Error(res, "Username can only have alphanumeric characters", http.StatusInternalServerError)
			Warning.Println("Username input is either empty or consists of illegal characters. Input: ", username)
			return
		}

		if alreadyInSession(username) {
			//User is already in session
			http.Error(res, "User is already in session.", http.StatusInternalServerError)
			Warning.Println("Already in session. Unable to start another one.")
			return
		}

		password := html.EscapeString(req.FormValue("password_input")) //Escape special characters

		var myUser User
		var ok error
		mutex.Lock() //Lock for global mapUsers read
		{
			db := OpenDB()
			defer db.Close()
			myUser, ok = GetUser(db, username) //Check if username exists
		}
		mutex.Unlock()

		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password)) //Compare password entered and password stored
		if ok != nil {
			http.Error(res, "No such user.", http.StatusForbidden)
			Warning.Println("No such user found: ", username)
			return
		} else if err != nil {
			http.Error(res, "Username and/or password do not match", http.StatusForbidden)
			Warning.Println("Username and/or password do not match. Username: ", username)
			return
		}

		mutex.Lock() //Lock for global mapSessions read
		{
			db := OpenDB()
			defer db.Close()
			_, err = GetSession(db, username)
		}
		mutex.Unlock()

		if err == nil {
			http.Error(res, "User already logged in.", http.StatusForbidden)
			Warning.Println("User:", username, "already logged in.")
			return
		}

		//Create session
		id, _ := uuid.NewV4() //Create Session ID
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  time.Now().AddDate(0, 0, 1), //Cookie will only last for a day
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}
		http.SetCookie(res, myCookie) //Set Cookie

		mutex.Lock() //Lock for global mapUsers read
		{
			db := OpenDB()
			defer db.Close()
			InsertSession(db, myCookie.Value, username) //Add the cookie object into Sesssions map, Key is the Session ID, Value is the username
		}
		mutex.Unlock()

		//Redirect to index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Info.Println("Successfully logged in. Username: ", username)
		return
	}

	tpl.ExecuteTemplate(res, "login.gohtml", nil) //Template: login.gohtml
}

func signupHTML(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		//User is already signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Already logged in. Unable to log in again.")
		return
	}

	var myUser User
	//Process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		if x, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username); !x || username == "" { //Regexp: Alphanumeric
			http.Error(res, "Username can only have alphanumeric characters", http.StatusInternalServerError)
			Warning.Println("Username input is either empty or consists of illegal characters. Input: ", username)
			return
		}

		password := html.EscapeString(req.FormValue("password")) //Escape special characters

		firstname := strings.Title(req.FormValue("firstname"))
		if x, _ := regexp.MatchString("^[a-zA-Z]+$", firstname); !x || firstname == "" { //Regexp: Alphabetical
			http.Error(res, "First name can only have english letters", http.StatusInternalServerError)
			Warning.Println("First name input is either empty or consists of illegal characters.")
			return
		}

		lastname := strings.Title(req.FormValue("lastname"))
		if x, _ := regexp.MatchString("^[a-zA-Z]+$", lastname); !x || lastname == "" { //Regexp: Alphabetical
			http.Error(res, "Last name can only have english letters", http.StatusInternalServerError)
			Warning.Println("Last name input is either empty or consists of illegal characters.")
			return
		}

		email := strings.Title(req.FormValue("email"))
		if x, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email); !x || email == "" { //Regexp: Email
			http.Error(res, "Please enter a valid email address.", http.StatusInternalServerError)
			Warning.Println("Email input is either not in the correct structure or is empty.")
			return
		}

		var ok error
		var db *sql.DB
		mutex.Lock() //Lock for global mapUsers read
		{
			db = OpenDB()
			defer db.Close()
			_, ok = GetUser(db, username) //Check if username exists
		}
		mutex.Unlock()

		if ok == nil { //Check if username exist/taken
			http.Error(res, "Username already taken", http.StatusForbidden)
			Warning.Println("Username ", username, " is already taken.")
			return
		}

		bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			Error.Println("Internal server error: ", err)
			return
		}
		InsertUser(db, username, bPassword, firstname, lastname, email)

		//Create session
		id, _ := uuid.NewV4() //Create Session ID
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  time.Now().AddDate(0, 0, 1), //Cookie will only last for a day
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}
		http.SetCookie(res, myCookie)               //Set Cookie
		InsertSession(db, myCookie.Value, username) //Add the cookie object into Sesssions map, Key is the Session ID, Value is the username

		//Redirect to index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(res, "signup.gohtml", myUser) //Template: signup.gohtml
}

func logoutHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	myCookie, _ := req.Cookie("myCookie")

	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		DeleteSession(db, myCookie.Value) //Delete the session from db
	}
	mutex.Unlock()

	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie) //Delete session cookie
	Info.Println("Successfully logged out. User session cookies have been deleted.")

	//Redirect to index
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func getUser(res http.ResponseWriter, req *http.Request) User { //Get current session cookie
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id, _ := uuid.NewV4() //Create Session ID
		myCookie = &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  time.Now().AddDate(0, 0, 1), //Cookie will only last for a day
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}
		http.SetCookie(res, myCookie)
		Info.Println("New session cookie created.")
	}

	var sessionx Session
	var userx User
	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		sessionx, err = GetSession(db, myCookie.Value)
		userx, err = GetUser(db, sessionx.Username) //Check if username exists
	}
	mutex.Unlock()

	return userx
}

//Check if user is already logged in
func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		Warning.Println("Internal server error: ", err)
		return false
	}

	var ok error
	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		_, ok = GetSession(db, myCookie.Value)
	}
	mutex.Unlock()

	if ok == nil { //Check if username exists
		Info.Println("User logged in: ", true)
		return true
	}

	Info.Println("User logged in: ", false)
	return false
}

//Check if user is already in session
func alreadyInSession(U string) bool {
	var err error
	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()

		var sessionx Session
		results := db.QueryRow(`SELECT * FROM Sessions WHERE Username = ?`, U)
		err = results.Scan(&sessionx.ID, &sessionx.Username)
	}
	mutex.Unlock()

	switch err {
	case sql.ErrNoRows:
		Warning.Println(sql.ErrNoRows)
	case nil:
		Info.Println("User already in session: ", true)
		return true
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return false
}

func accountHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	tpl.ExecuteTemplate(res, "account.gohtml", nil)
}

func editUserHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	userx := getUser(res, req)

	if req.Method == http.MethodPost { // get form values
		var err error
		bPassword := userx.Password
		if req.FormValue("password") != "" {
			bPassword, err = bcrypt.GenerateFromPassword([]byte(html.EscapeString(req.FormValue("password"))), bcrypt.MinCost)
			if bPassword != nil && err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				Error.Println("Internal server error: ", err)
				return
			}
		}

		firstname := strings.Title(req.FormValue("firstname"))
		if x, _ := regexp.MatchString("^[a-zA-Z]+$", firstname); !x || firstname == "" { //Regexp: Alphabetical
			http.Error(res, "First name can only have english letters", http.StatusInternalServerError)
			Warning.Println("First name input is either empty or consists of illegal characters.")
			return
		}

		lastname := strings.Title(req.FormValue("lastname"))
		if x, _ := regexp.MatchString("^[a-zA-Z]+$", lastname); !x || lastname == "" { //Regexp: Alphabetical
			http.Error(res, "Last name can only have english letters", http.StatusInternalServerError)
			Warning.Println("Last name input is either empty or consists of illegal characters.")
			return
		}

		email := strings.Title(req.FormValue("email"))
		if x, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email); !x || email == "" { //Regexp: Email
			http.Error(res, "Please enter a valid email address.", http.StatusInternalServerError)
			Warning.Println("Email input is either not in the correct structure or is empty.")
			return
		}

		var ok error
		var db *sql.DB
		mutex.Lock() //Lock for global mapUsers read
		{
			db = OpenDB()
			defer db.Close()
			_, ok = GetUserEmail(db, email)          //Check if email exists
			if ok == nil && (email != userx.Email) { //Check if email exist/taken
				http.Error(res, "Email already in use", http.StatusForbidden)
				Warning.Println("Email ", email, " is in use.")
				return
			}

			UpdateUser(db, userx.Username, bPassword, firstname, lastname, email)
		}
		mutex.Unlock()

		http.Redirect(res, req, "/account", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "editUser.gohtml", userx)
}

func delUserHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { //Check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	username := req.URL.Query().Get("id")
	myCookie, _ := req.Cookie("myCookie")
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		DeleteSession(db, myCookie.Value) //Delete the session from db
		DeleteUser(db, username) //Delete the user from db
	}
	mutex.Unlock()

	http.Redirect(res, req, "/", http.StatusSeeOther)
}