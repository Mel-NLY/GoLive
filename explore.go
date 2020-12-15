package main

import (
	"net/http"

	"GoLive/pkgs"
)

func exploreHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var postsx pkgs.BST
	var tempPosts []pkgs.Post

	userx := getUser(res, req) //Get user

	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		postsx = GetPosts(db)
	}
	mutex.Unlock()

	tempPosts = postsx.InOrder()

	TempStruct := struct {
		Userx  User
		Postsx []pkgs.Post
	}{
		userx,
		tempPosts,
	}

	tpl.ExecuteTemplate(res, "explore.gohtml", TempStruct)
}
