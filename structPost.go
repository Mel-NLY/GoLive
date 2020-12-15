package main

import (
	"GoLive/pkgs"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/xconstruct/go-pushbullet"
	gomail "gopkg.in/mail.v2"
)

func viewPostHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var postx pkgs.Post
	var err error
	id := req.URL.Query().Get("id")

	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		postx, err = GetPost(db, id)
	}
	mutex.Unlock()

	if (err != nil && postx != pkgs.Post{}) {
		http.Error(res, "Invalid Post ID", http.StatusForbidden)
		Warning.Println("PostID:", id, "is invalid.")
		return
	}

	tpl.ExecuteTemplate(res, "viewPost.gohtml", postx)
}

func createPostHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	tpl.ExecuteTemplate(res, "createPost.gohtml", nil)
}

func uploadFileHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	Warning.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	req.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := req.FormFile("myFile")
	if err != nil {
		Warning.Println("Error Retrieving the File")
		return
	}
	defer file.Close()
	Info.Println("Uploaded File:", handler.Filename)
	Info.Println("File Size:", handler.Size)
	Info.Println("MIME Header:", handler.Header)

	// Create a temporary file within our posts directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("templates/assets/img/posts", "upload-*.png")
	if err != nil {
		Warning.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		Warning.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	Info.Println("Successfully Uploaded File.")

	userx := getUser(res, req)

	title := req.FormValue("title")
	if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", title); !x || title == "" { //Regexp: Alphanumeric
		http.Error(res, "Title consists of illegal characters", http.StatusInternalServerError)
		Warning.Println("Title input is either empty or consists of illegal characters. Input: ", title)
		return
	}

	fn := tempFile.Name()[26:]

	des := req.FormValue("description")
	if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", des); !x || des == "" { //Regexp: Alphanumeric
		http.Error(res, "Description consists of illegal characters", http.StatusInternalServerError)
		Warning.Println("Description input is either empty or consists of illegal characters. Input: ", des)
		return
	}

	tag := req.FormValue("tag")

	mutex.Lock() //Lock for global mapUsers read
	{
		//Send data to SQL
		db := OpenDB()
		defer db.Close()
		InsertPost(db, userx.Username, title, fn, des, tag)
	}
	mutex.Unlock()

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func myPostsHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var postsx pkgs.BST

	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		postsx = GetPosts(db)
	}
	mutex.Unlock()

	userx := getUser(res, req)

	tempPosts := postsx.PreOrder(userx.Username)

	tpl.ExecuteTemplate(res, "myPosts.gohtml", tempPosts)
}

func contactHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var userx User
	id := req.URL.Query().Get("id")
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		userx, _ = GetUser(db, id)
	}
	mutex.Unlock()

	if req.Method == http.MethodPost { // get form values
		message := req.FormValue("message")
		if message == "" {
			http.Error(res, "Message body is empty", http.StatusInternalServerError)
			Warning.Println("Message body is empty.")
			return
		}
		phoneno := req.FormValue("phoneno")
		pushBullet(getUser(res, req), message, phoneno) //Send pushbullet message to admin

		sendEmail(getUser(res, req), userx, message, phoneno) //Send email to user

		http.Redirect(res, req, "/explore", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "contact.gohtml", userx)
}

func pushBullet(userx User, message string, phoneno string) {
	pb := pushbullet.New("o.KORUKnutV6FhKO3uEaYFkZY4Tm12jB1S")
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	} else {
		err = pb.PushNote(devs[0].Iden, "New message from: "+userx.Username, message+phoneno)
		if err != nil {
			panic(err)
		}
	}
}

func editPostHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var postx pkgs.Post
	var err error
	id := req.URL.Query().Get("id")

	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		postx, err = GetPost(db, id)
	}
	mutex.Unlock()

	if (err != nil && postx != pkgs.Post{}) {
		http.Error(res, "Invalid Post ID", http.StatusForbidden)
		Warning.Println("PostID:", id, "is invalid.")
		return
	}

	if req.Method == http.MethodPost { // get form values
		/*UPDATE SQL DATA ++ MAKE EMAIL A UNIQUE KEY*/
		fn := postx.Image
		if req.FormValue("myFileName") != "" { //If value is empty, image is unchanged
			Warning.Println("File Upload Endpoint Hit")

			// Parse our multipart form, 10 << 20 specifies a maximum
			// upload of 10 MB files.
			req.ParseMultipartForm(10 << 20)
			// FormFile returns the first file for the given key `myFile`
			// it also returns the FileHeader so we can get the Filename,
			// the Header and the size of the file
			file, handler, err := req.FormFile("myFile")
			if err != nil {
				Warning.Println("Error Retrieving the File")
				log.Fatalln("Failed to open file:", err)
			}
			defer file.Close()
			Info.Println("Uploaded File:", handler.Filename)
			Info.Println("File Size:", handler.Size)
			Info.Println("MIME Header:", handler.Header)

			// Create a temporary file within our posts directory that follows
			// a particular naming pattern
			tempFile, err := ioutil.TempFile("templates/assets/img/posts", "upload-*.png")
			if err != nil {
				Error.Println(err)
				log.Fatalln("Failed to create file:", err)
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				Error.Println(err)
				log.Fatalln("Failed to read file:", err)
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)
			// return that we have successfully uploaded our file!
			Info.Println("Successfully Uploaded File.")

			fn = tempFile.Name()[26:]
		}

		userx := getUser(res, req)

		title := req.FormValue("title")
		if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", title); !x || title == "" { //Regexp: Alphanumeric
			http.Error(res, "Title consists of illegal characters", http.StatusInternalServerError)
			Warning.Println("Title input is either empty or consists of illegal characters. Input: ", title)
			return
		}

		des := req.FormValue("description")
		if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", des); !x || des == "" { //Regexp: Alphanumeric
			http.Error(res, "Description consists of illegal characters", http.StatusInternalServerError)
			Warning.Println("Description input is either empty or consists of illegal characters. Input: ", des)
			return
		}

		var tag string
		if req.FormValue("tag") != "" {
			tag = req.FormValue("tag")
		} else {
			tag = postx.Tag
		}

		mutex.Lock() //Lock for global mapUsers read
		{
			//Send data to SQL
			db := OpenDB()
			defer db.Close()
			UpdatePost(db, postx.ID, userx.Username, title, fn, des, tag)
		}
		mutex.Unlock()

		http.Redirect(res, req, "/", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "editPost.gohtml", postx)
}

func sendEmail(userFrom User, userTo User, message string, phoneno string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "biketransport.bt@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "12.melissa.2i2.pw14@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Message from "+userFrom.Username)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", message)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "biketransport.bt@gmail.com", "Biketransport#1")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}
