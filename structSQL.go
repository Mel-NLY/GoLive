package main

import (
	"GoLive/pkgs"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Session is exported
type Session struct {
	ID       []uint8
	Username string
}

//GetUsers retrieves all records in db
func GetUsers(db *sql.DB) {
	results, err := db.Query("SELECT * FROM BikeTransport_db.Users")
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	for results.Next() {
		var userx User
		err := results.Scan(&userx.Username, &userx.Password, &userx.First, &userx.Last, &userx.Email)
		if err != nil {
			Error.Println("Failed to scan in SQL records:", err)
			log.Fatalln("Failed to scan in SQL records:", err)
		}
	}
	Info.Println("User records successfully retrieved.")
}

//GetUser retrieves one specific record in db
func GetUser(db *sql.DB, U string) (User, error) {
	var userx User
	results := db.QueryRow(`SELECT * FROM Users WHERE Username=?`, U)
	err := results.Scan(&userx.Username, &userx.Password, &userx.First, &userx.Last, &userx.Email)
	switch err {
	case sql.ErrNoRows:
		return User{}, errors.New("No User records found for Username = " + U)
	case nil:
		return userx, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return User{}, err
}

//GetUserEmail retrieves one specific record in db
func GetUserEmail(db *sql.DB, E string) (User, error) {
	var userx User
	results := db.QueryRow(`SELECT * FROM Users WHERE Email=?`, E)
	err := results.Scan(&userx.Username, &userx.Password, &userx.First, &userx.Last, &userx.Email)
	switch err {
	case sql.ErrNoRows:
		return User{}, errors.New("No User records found for Email = " + E)
	case nil:
		return userx, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return User{}, err
}

//InsertUser inserts the specified record to the db
func InsertUser(db *sql.DB, U string, P []byte, F string, L string, E string) {
	query := fmt.Sprintf("INSERT INTO Users VALUES ('%s', '%s', '%s', '%s', '%s')", U, string(P), F, L, E)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("User - ", U, "successfully inserted.")
}

//UpdateUser edit specified record to the db
func UpdateUser(db *sql.DB, U string, P []byte, F string, L string, E string) {
	query := fmt.Sprintf("UPDATE Users SET Password='%s', First='%s', Last='%s', Email='%s' WHERE Username='%s'", P, F, L, E, U)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("User - ", U, "successfully updated.")
}

//DeleteUser deletes records of the specified user id
// func DeleteUser(db *sql.DB, U string) {
// 	query := fmt.Sprintf("DELETE FROM Users WHERE Username='%s'", U)
// 	_, err := db.Query(query)
// 	if err != nil {
// 		Error.Println("Failed to execute SQL command:", err)
// 		log.Fatalln("Failed to execute SQL command:", err)
// 	}
// 	Info.Println("User - ", U, "successfully deleted.")
// }

//GetSession checks if session key exists in the db
func GetSession(db *sql.DB, id string) (Session, error) {
	var sessionx Session
	results := db.QueryRow(`SELECT * FROM Sessions WHERE SessionID = ?`, id)
	err := results.Scan(&sessionx.ID, &sessionx.Username)
	switch err {
	case sql.ErrNoRows:
		return Session{}, errors.New("No Session records found for ID = " + id)
	case nil:
		return sessionx, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return Session{}, err
}

//InsertSession inserts the specified record to the db
func InsertSession(db *sql.DB, id string, U string) {
	query := fmt.Sprintf("INSERT INTO Sessions VALUES ('%s', (SELECT Username from Users WHERE Username='%s'))", id, U)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Session - ", id, "successfully inserted.")
}

//DeleteSession deletes records of the specified user id
func DeleteSession(db *sql.DB, id string) {
	query := fmt.Sprintf("DELETE FROM Sessions WHERE SessionID='%s'", id)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Session - ", id, "successfully deleted.")
}

//GetPosts retrieves all records in db
func GetPosts(db *sql.DB) pkgs.BST {
	results, err := db.Query("SELECT * FROM BikeTransport_db.Posts")
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	var postsx pkgs.BST
	for results.Next() {
		var postx pkgs.Post
		err := results.Scan(&postx.ID, &postx.Username, &postx.Time.Day, &postx.Time.Month, &postx.Time.Year, &postx.Time.Hour, &postx.Time.Min, &postx.Title, &postx.Image, &postx.Description, &postx.Tag)
		if err != nil {
			Error.Println("Failed to scan in SQL records:", err)
			log.Fatalln("Failed to scan in SQL records:", err)
		}
		postsx.Insert(postx)
	}
	Info.Println("Post records successfully retrieved.")
	return postsx
}

//GetPost retrieves one specific record in db
func GetPost(db *sql.DB, id string) (pkgs.Post, error) {
	var postx pkgs.Post
	results := db.QueryRow(`SELECT * FROM Posts WHERE PostID=?`, id)
	err := results.Scan(&postx.ID, &postx.Username, &postx.Time.Day, &postx.Time.Month, &postx.Time.Year, &postx.Time.Hour, &postx.Time.Min, &postx.Title, &postx.Image, &postx.Description, &postx.Tag)
	switch err {
	case sql.ErrNoRows:
		return pkgs.Post{}, errors.New("No Post records found for ID = " + id)
	case nil:
		return postx, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return pkgs.Post{}, err
}

//InsertPost inserts the specified record to the db
func InsertPost(db *sql.DB, U string, T string, I string, D string, Tag string) {
	// Example: this will give us a 44 byte, base64 encoded output
	id, err := GenerateRandomString(32)
	if err != nil {
		Warning.Println(err)
		return
	}

	query := fmt.Sprintf("INSERT INTO Posts VALUES ('%s', '%s', %d, %d, %d, %d, %d, '%s', '%s', '%s', '%s')", id, U, time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), T, I, D, Tag)
	_, err = db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Post ID - ", id, "successfully inserted.")
}

//UpdatePost inserts the specified record to the db
func UpdatePost(db *sql.DB, id string, U string, T string, I string, D string, Tag string) {

	query := fmt.Sprintf("UPDATE Posts SET Username='%s', TimeDay=%d, TimeMonth=%d, TimeYear=%d, TimeHour=%d, TimeMin=%d, Title='%s', Image='%s', Description='%s', Tag='%s' WHERE PostID='%s'", U, time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), T, I, D, Tag, id)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Post ID - ", id, "successfully updated.")
}

//DeletePost deletes records of the specified post id
func DeletePost(db *sql.DB, id string) {
	query := fmt.Sprintf("DELETE FROM Posts WHERE PostID='%s'", id)
	_, err := db.Query(query)
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Post - ", id, "successfully deleted.")
}

//GetRoutes retrieves all records in db
func GetRoutes(db *sql.DB) []Route {
	results, err := db.Query("SELECT * FROM BikeTransport_db.Routes ORDER BY RouteID")
	if err != nil {
		Error.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	var routesx []Route
	for results.Next() {
		var routex Route
		err := results.Scan(&routex.ID, &routex.Username, &routex.Duration, &routex.Distance, &routex.Name)
		if err != nil {
			Error.Println("Failed to scan in SQL records:", err)
			log.Fatalln("Failed to scan in SQL records:", err)
		}
		routesx = append(routesx, routex)
	}
	Info.Println("Post records successfully retrieved.")
	return routesx
}

//GetRoute retrieves one specific record in db
func GetRoute(db *sql.DB, id string) (Route, error) {
	var routex Route
	results := db.QueryRow(`SELECT * FROM Routes WHERE RouteID=?`, id)
	err := results.Scan(&routex.ID, &routex.Username, &routex.Distance, &routex.Duration, &routex.Name)
	switch err {
	case sql.ErrNoRows:
		return Route{}, errors.New("No Route records found for ID = " + id)
	case nil:
		return routex, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}

	return Route{}, err
}

//InsertRoute inserts the specified record to the db
func InsertRoute(db *sql.DB, id string, U string) {
	query := fmt.Sprintf("INSERT INTO Routes VALUES ('%s', '%s', NULL, NULL, NULL)", id, U)
	_, err := db.Query(query)
	if err != nil {
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Route ID - ", id, "successfully inserted.")
}

//UpdateRoute inserts the specified record to the db
func UpdateRoute(db *sql.DB, id string, U string, DS float64, D float64, RN string) {
	query := fmt.Sprintf("UPDATE Routes SET Username='%s', DurationSec=%f, Distance=%f, RouteName='%s' WHERE RouteID='%s'", U, DS, D, RN, id)
	_, err := db.Query(query)
	if err != nil {
		Info.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Route ID - ", id, "successfully updated.")
}

//DeleteRoute inserts the specified record to the db
func DeleteRoute(db *sql.DB, id string) {
	query := fmt.Sprintf("DELETE FROM RoutePoints WHERE RouteID='%s'", id)
	_, err := db.Query(query)
	if err != nil {
		Info.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	query = fmt.Sprintf("DELETE FROM Routes WHERE RouteID='%s'", id)
	_, err = db.Query(query)
	if err != nil {
		Info.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	Info.Println("Route ID - ", id, "and RoutePoints successfully deleted.")
}

//GetRoutePoints retrieves all records in db
func GetRoutePoints(db *sql.DB, routeid string) ([]pkgs.RoutePoint, float64) {
	results, err := db.Query("SELECT * FROM BikeTransport_db.RoutePoints WHERE RouteID=? ORDER BY Position", routeid)
	if err != nil {
		Info.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	var distance float64
	var route pkgs.LinkedList
	var routePx pkgs.RoutePoint
	var listRP []pkgs.RoutePoint
	for results.Next() {
		err := results.Scan(&routePx.ID, &routePx.RouteID, &routePx.Lat, &routePx.Lon, &routePx.Time.Day, &routePx.Time.Month, &routePx.Time.Year, &routePx.Time.Hour, &routePx.Time.Min, &routePx.Position)
		if err != nil {
			Info.Println("Failed to scan in SQL records:", err)
			log.Fatalln("Failed to scan in SQL records:", err)
		}
		route.AddNode(routePx)
	}

	currentNode := route.Head
	nextNode := currentNode.Next
	for nextNode.Next != nil {
		Info.Println(currentNode.Item.Lat, currentNode.Item.Lon)
		Info.Println(nextNode.Item.Lat, nextNode.Item.Lon)
		distance += findDistance(currentNode.Item.Lat, currentNode.Item.Lon, nextNode.Item.Lat, nextNode.Item.Lon)
		listRP = append(listRP, currentNode.Item)
		currentNode = currentNode.Next
		nextNode = nextNode.Next
	}
	Info.Println("Distance calculated: ", distance)
	return listRP, distance
}

//GetRoutePoint retrieves one specific record in db
func GetRoutePoint(db *sql.DB, id string, waypoints []int) (pkgs.RoutePoint, pkgs.RoutePoint, pkgs.LinkedList, error) {
	var routeStartPx pkgs.RoutePoint
	var routeEndPx pkgs.RoutePoint
	var routeWayPx pkgs.LinkedList

	results := db.QueryRow(`SELECT * FROM RoutePoints WHERE RouteID=? AND Position=1;`, id)
	err := results.Scan(&routeStartPx.ID, &routeStartPx.RouteID, &routeStartPx.Lat, &routeStartPx.Lon, &routeStartPx.Time.Day, &routeStartPx.Time.Month, &routeStartPx.Time.Year, &routeStartPx.Time.Hour, &routeStartPx.Time.Min, &routeStartPx.Position)

	results = db.QueryRow(`SELECT *	FROM RoutePoints WHERE RouteID=? AND Position=(SELECT MAX(Position) FROM RoutePoints WHERE RouteID=?);`, id, id)
	err = results.Scan(&routeEndPx.ID, &routeEndPx.RouteID, &routeEndPx.Lat, &routeEndPx.Lon, &routeEndPx.Time.Day, &routeEndPx.Time.Month, &routeEndPx.Time.Year, &routeEndPx.Time.Hour, &routeEndPx.Time.Min, &routeEndPx.Position)

	for _, c := range waypoints {
		var routePx pkgs.RoutePoint
		results = db.QueryRow(`SELECT * FROM RoutePoints WHERE RouteID=? AND Position=?;`, id, c)
		err = results.Scan(&routePx.ID, &routePx.RouteID, &routePx.Lat, &routePx.Lon, &routePx.Time.Day, &routePx.Time.Month, &routePx.Time.Year, &routePx.Time.Hour, &routePx.Time.Min, &routePx.Position)
		routeWayPx.AddNode(routePx)
	}

	switch err {
	case sql.ErrNoRows:
		return pkgs.RoutePoint{}, pkgs.RoutePoint{}, pkgs.LinkedList{}, errors.New("No Route records found for ID = " + id)
	case nil:
		return routeStartPx, routeEndPx, routeWayPx, nil
	default:
		Error.Println(err)
		log.Fatal(err)
	}
	return pkgs.RoutePoint{}, pkgs.RoutePoint{}, pkgs.LinkedList{}, err
}

//InsertRoutePoint inserts the specified record to the db
func InsertRoutePoint(db *sql.DB, routeid string, Lat float64, Lon float64, P int) {
	id, err := GenerateRandomString(32)
	if err != nil {
		Info.Println(err)
		return
	}

	query := fmt.Sprintf("INSERT INTO RoutePoints VALUES ('%s', '%s', %.10f, %.10f, %d, %d, %d, %d, %d, %d)", id, routeid, Lat, Lon, time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), P)
	_, err = db.Query(query)
	if err != nil {
		Info.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	Info.Println("Point for Route ID - ", routeid, "successfully inserted.")
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//OpenDB would open the db but the close has to be done after the action
func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(192.168.1.200:32769)/BikeTransport_db")
	if err != nil {
		Error.Println("Failed to open SQL db:", err)
		log.Fatalln("Failed to open SQL db:", err)
	}
	Info.Println("Database successfully opened.")

	return db
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func findDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
