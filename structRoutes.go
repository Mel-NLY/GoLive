package main

import (
	"GoLive/pkgs"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Route is exported
type Route struct {
	ID       string
	Name     sql.NullString
	Username string
	Duration sql.NullFloat64
	Distance sql.NullFloat64
}

var routeid string
var stopTracker bool
var timerChannel chan bool
var duration time.Duration
var position int

func routesHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	timerChannel = make(chan bool) //Creation of timerChannel for structRoutes

	var routesx []Route
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		routesx = GetRoutes(db)
	}
	mutex.Unlock()

	tempStruct := struct{
		Routesx []Route
		Userx User
	}{
		routesx,
		getUser(res, req),
	}

	tpl.ExecuteTemplate(res, "routes.gohtml", tempStruct)
}

func repeatedHome(res http.ResponseWriter, req *http.Request) {
	go func() {
		timeStart := time.Now()

		<-timerChannel

		timeEnd := time.Now()
		duration = timeEnd.Sub(timeStart)
		return
	}()
	position = -1
	routeid = ""
	http.Redirect(res, req, "/tracking", http.StatusSeeOther)
	return
}

func home(res http.ResponseWriter, req *http.Request) {
	position++
	tpl.ExecuteTemplate(res, "getLocation.gohtml", nil)
}

func location(res http.ResponseWriter, req *http.Request) {
	var err error
	vars := mux.Vars(req)
	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	long, _ := strconv.ParseFloat(vars["long"], 64)

	if routeid == "" {
		// Example: this will give us a 44 byte, base64 encoded output
		routeid, err = GenerateRandomString(32)
		if err != nil {
			Warning.Println(err)
			return
		}

		userx := getUser(res, req)
		mutex.Lock()
		{
			db := OpenDB()
			defer db.Close()
			InsertRoute(db, routeid, userx.Username) //Creating new route record
		}
		mutex.Unlock()
	}

	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		InsertRoutePoint(db, routeid, lat, long, position) //Adding new route point record every 10sec
	}
	mutex.Unlock()

	time.Sleep(5 * time.Second)

	http.Redirect(res, req, "/tracking", http.StatusSeeOther)
	return
}

func viewRouteHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var routeStartPx pkgs.RoutePoint
	var routeEndPx pkgs.RoutePoint
	var routeWayPx pkgs.LinkedList

	var count float64
	var max float64
	var points float64
	var waypoints []int
	var err error
	id := req.URL.Query().Get("id")

	mutex.Lock() //Lock for global mapUsers read
	{
		db := OpenDB()
		defer db.Close()
		results := db.QueryRow(`SELECT MAX(Position) FROM RoutePoints WHERE RouteID=?;`, id)
		err = results.Scan(&count)
		max = math.Round(count / 20)
		for count > points && max > 1 { //Obtain waypoints
			waypoints = append(waypoints, int(points))
			points += max
		}
		sort.Ints(waypoints)
		routeStartPx, routeEndPx, routeWayPx, _ = GetRoutePoint(db, id, waypoints)
	}
	mutex.Unlock()

	var wayPt string
	currentNode := routeWayPx.Head
	if currentNode != nil { //Obtain waypoint lat long
		wayPt = fmt.Sprintf("%f,%f", currentNode.Next.Item.Lat, currentNode.Next.Item.Lon)
		for currentNode.Next != nil {
			temp := fmt.Sprintf("|%f,%f", currentNode.Next.Item.Lat, currentNode.Next.Item.Lon)
			wayPt = fmt.Sprintf("%s%s", wayPt, temp)
			currentNode = currentNode.Next
		}
	}

	if err != nil {
		http.Error(res, "Invalid Route ID", http.StatusForbidden)
		Warning.Println("RouteID:", id, "is invalid.")
		return
	}

	tempStruct := struct {
		RouteStartPx pkgs.RoutePoint
		RouteEndPx   pkgs.RoutePoint
		RouteWayPx   string
	}{
		routeStartPx,
		routeEndPx,
		wayPt,
	}

	tpl.ExecuteTemplate(res, "viewRoute.gohtml", tempStruct)
}

func createRouteHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	if req.Method == http.MethodPost { // get form values
		routeN := req.FormValue("routeName")
		if x, _ := regexp.MatchString("^(.|\\s)*[a-zA-Z]+(.|\\s)*$", routeN); !x || routeN == "" { //Regexp: Alphanumeric
			http.Error(res, "Route name consists of illegal characters", http.StatusInternalServerError)
			Warning.Println("Route name input is either empty or consists of illegal characters. Input: ", routeN)
			return
		}

		userx := getUser(res, req)
		mutex.Lock()
		{
			db := OpenDB()
			defer db.Close()
			_, distance := GetRoutePoints(db, routeid)
			UpdateRoute(db, routeid, userx.Username, duration.Seconds(), distance, routeN)
		}
		mutex.Unlock()

		http.Redirect(res, req, "/routes", http.StatusSeeOther)
		return
	}

	timerChannel <- false //Stop tracking

	tpl.ExecuteTemplate(res, "createRoute.gohtml", nil)
}

func delRouteHTML(res http.ResponseWriter, req *http.Request) {
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
		DeleteRoute(db, id)
	}
	mutex.Unlock()

	http.Redirect(res, req, "/routes", http.StatusSeeOther)
}

func heatmapHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	var listR []Route
	var listRP []pkgs.RoutePoint
	mutex.Lock()
	{
		db := OpenDB()
		defer db.Close()
		listR = GetRoutes(db)
		for _, c := range listR {
			x, _ := GetRoutePoints(db, c.ID)
			listRP = append(listRP, x...)
		}
	}
	mutex.Unlock()

	tpl.ExecuteTemplate(res, "heatmap.gohtml", listRP)
}
