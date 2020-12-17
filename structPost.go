package main

import (
	"GoLive/pkgs"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

	tempStruct := struct {
		Postx pkgs.Post
		Userx User
	}{
		postx,
		getUser(res, req),
	}

	tpl.ExecuteTemplate(res, "viewPost.gohtml", tempStruct)
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

	var fn string
	if req.FormValue("myFileName") == "" { //Check if value is empty
		http.Error(res, "Please upload an image file", http.StatusInternalServerError)
		Warning.Println("No file uploaded.")
		return
	}

	fn = uploadFile(res, req)
	userx := getUser(res, req)
	title, des, err := checkTitleDes(res, req)
	if err != nil{
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

func uploadFile(res http.ResponseWriter, req *http.Request) string {
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
		log.Fatalln("Failed to create file:", err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	Info.Println("Successfully Uploaded File.")
	fn := tempFile.Name()[26:]

	return fn
}

func checkTitleDes(res http.ResponseWriter, req *http.Request) (string, string, error) {
	title := req.FormValue("title")
	if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", title); !x || title == "" { //Regexp: Alphanumeric
		http.Error(res, "Title consists of illegal characters or is empty", http.StatusInternalServerError)
		Warning.Println("Title input is either empty or consists of illegal characters. Input: ", title)
		return "", "", errors.New("Title consists of illegal characters or is empty")
	}
	title = strings.Replace(title, "'", "\\'", -1) //Escape single quotes

	des := req.FormValue("description")
	if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", des); !x || des == "" { //Regexp: Alphanumeric
		http.Error(res, "Description consists of illegal characters or is empty", http.StatusInternalServerError)
		Warning.Println("Description input is either empty or consists of illegal characters. Input: ", des)
		return "", "", errors.New("Description consists of illegal characters or is empty")
	}
	des = strings.Replace(des, "'", "\\'", -1) //Escape single quotes

	return title, des, nil
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
	var postx pkgs.Post
	username := req.URL.Query().Get("user")
	postid := req.URL.Query().Get("post")
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		userx, _ = GetUser(db, username)
		postx, _ = GetPost(db, postid)
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
		pushBullet(getUser(res, req), userx, postx, message, phoneno) //Send pushbullet message to admin

		sendEmail(getUser(res, req), userx, postx, message, phoneno) //Send email to user

		http.Redirect(res, req, "/explore", http.StatusSeeOther)
	}

	tpl.ExecuteTemplate(res, "contact.gohtml", userx)
}

func pushBullet(userFrom User, userTo User, postAbout pkgs.Post, message string, phoneno string) {
	pb := pushbullet.New("o.KORUKnutV6FhKO3uEaYFkZY4Tm12jB1S")
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	} else {
		body := fmt.Sprintln("Hi there,\n\nA BikeTransport user has sent you the following message:\n\nPost Title: " + postAbout.Title + "\nPost Description: " + postAbout.Description + "\nPost Created/Edited on: " + strconv.Itoa(postAbout.Time.Day) + "-" + strconv.Itoa(postAbout.Time.Month) + "-" + strconv.Itoa(postAbout.Time.Year) + " " + strconv.Itoa(postAbout.Time.Hour) + ":" + strconv.Itoa(postAbout.Time.Min) + "\nPost Tag: " + postAbout.Tag + "\n\nMessage from " + userFrom.Username + ":\n" + message + "\nEmail: " + userFrom.Email + "\nPhone Number(Optional): " + phoneno + "\n\n\n\nHave a great day,\nThe BikeTransport Team\nPlease do not reply to this message.")
		err = pb.PushNote(devs[0].Iden, userFrom.Username+" has received a new message from "+userTo.Username, body)
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

	if req.Method == http.MethodPost { // Get form values
		fn := postx.Image
		if req.FormValue("myFileName") != "" { //If value is empty, image is unchanged
			fn = uploadFile(res, req)
		}

		userx := getUser(res, req)

		title, des, err := checkTitleDes(res, req)
		if err != nil{
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

func delPostHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	id := req.URL.Query().Get("id")
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		DeletePost(db, id)
	}
	mutex.Unlock()

	http.Redirect(res, req, "/myPosts", http.StatusSeeOther)
}

func sendEmail(userFrom User, userTo User, postAbout pkgs.Post, message string, phoneno string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "biketransport.bt@gmail.com")                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            //Email Sender
	m.SetHeader("To", userTo.Email)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              //Email Receiver(s)
	m.SetHeader("Subject", "Message from "+userFrom.Username)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    //Email Subject
	m.SetBody("text/plain", "Hi there,\n\nA BikeTransport user has sent you the following message:\n\nPost Title: "+postAbout.Title+"\nPost Description: "+postAbout.Description+"\nPost Created/Edited on: "+strconv.Itoa(postAbout.Time.Day)+"-"+strconv.Itoa(postAbout.Time.Month)+"-"+strconv.Itoa(postAbout.Time.Year)+" "+strconv.Itoa(postAbout.Time.Hour)+":"+strconv.Itoa(postAbout.Time.Min)+"\nPost Tag: "+postAbout.Tag+"\n\nMessage from "+userFrom.Username+":\n"+message+"\nEmail: "+userFrom.Email+"\nPhone Number(Optional): "+phoneno+"\n\n\n\nHave a great day,\nThe BikeTransport Team\nPlease do not reply to this email.") //Email body

	d := gomail.NewDialer("smtp.gmail.com", 587, "biketransport.bt@gmail.com", "Biketransport#1") //Settings for the SMTP server

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} //Only necessary when the SSL/TLS certificate is not valid on the server
	if err := d.DialAndSend(m); err != nil {            //Sending the email
		Error.Println(err)
		log.Fatalln("Failed to send email:", err)
	}
	return
}
