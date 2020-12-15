package main

import (
	"bufio"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"GoLive/pkgs"

	_ "github.com/go-sql-driver/mysql"
)

var in = bufio.NewReader(os.Stdin)
var position int

func main() {
	f, _ := os.Open("dummyRoutes.txt")
	scanner := bufio.NewScanner(f)
	routeid, _ := GenerateRandomString(32)
	db := OpenDB()
	defer f.Close()
	defer db.Close()

	InsertRoute(db, routeid, "user3")
	for scanner.Scan() { //Scanning each line of the file
		s := strings.Split(scanner.Text(), ",")
		lon, _ := strconv.ParseFloat(s[0], 64)
		lat, _ := strconv.ParseFloat(s[1], 64)
		InsertRoutePoint(db, routeid, lat, lon, position)
		position++
	}
	distance := GetRoutePoints(db, routeid)
	UpdateRoute(db, routeid, "user3", distance/7, distance, "Bukit Timah >> Ngee Ann")
	return
}

//InsertRoute inserts the specified record to the db
func InsertRoute(db *sql.DB, id string, U string) {
	query := fmt.Sprintf("INSERT INTO Routes VALUES ('%s', '%s', NULL, NULL, NULL)", id, U)
	_, err := db.Query(query)
	if err != nil {
		log.Fatalln("Failed to execute SQL command:", err)
	}
	fmt.Println("Route ID - ", id, "successfully inserted.")
}

//InsertRoutePoint inserts the specified record to the db
func InsertRoutePoint(db *sql.DB, routeid string, Lat float64, Lon float64, P int) {
	id, err := GenerateRandomString(32)
	if err != nil {
		fmt.Println(err)
		return
	}

	query := fmt.Sprintf("INSERT INTO RoutePoints VALUES ('%s', '%s', %.10f, %.10f, %d, %d, %d, %d, %d, %d)", id, routeid, Lat, Lon, time.Now().Day(), time.Now().Month(), time.Now().Year(), time.Now().Hour(), time.Now().Minute(), P)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	fmt.Println("Point for Route ID - ", routeid, "successfully inserted.")
}

//UpdateRoute inserts the specified record to the db
func UpdateRoute(db *sql.DB, id string, U string, DS float64, D float64, RN string) {
	query := fmt.Sprintf("UPDATE Routes SET Username='%s', DurationSec=%f, Distance=%f, RouteName='%s' WHERE RouteID='%s'", U, DS, D, RN, id)
	_, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}
	fmt.Println("Route ID - ", id, "successfully updated.")
}

//GetRoutePoints retrieves all records in db
func GetRoutePoints(db *sql.DB, routeid string) float64 {
	results, err := db.Query("SELECT * FROM BikeTransport_db.RoutePoints WHERE RouteID=? ORDER BY Position", routeid)
	if err != nil {
		fmt.Println("Failed to execute SQL command:", err)
		log.Fatalln("Failed to execute SQL command:", err)
	}

	var distance float64
	var route pkgs.LinkedList
	var routePx pkgs.RoutePoint
	for results.Next() {
		err := results.Scan(&routePx.ID, &routePx.RouteID, &routePx.Lat, &routePx.Lon, &routePx.Time.Day, &routePx.Time.Month, &routePx.Time.Year, &routePx.Time.Hour, &routePx.Time.Min, &routePx.Position)
		if err != nil {
			fmt.Println("Failed to scan in SQL records:", err)
			log.Fatalln("Failed to scan in SQL records:", err)
		}
		route.AddNode(routePx)
	}

	currentNode := route.Head
	nextNode := currentNode.Next
	for nextNode.Next != nil {
		fmt.Println(currentNode.Item.Lat, currentNode.Item.Lon)
		fmt.Println(nextNode.Item.Lat, nextNode.Item.Lon)
		distance += findDistance(currentNode.Item.Lat, currentNode.Item.Lon, nextNode.Item.Lat, nextNode.Item.Lon)
		currentNode = currentNode.Next
		nextNode = nextNode.Next
	}
	fmt.Println("Distance calculate: ", distance)
	return distance
}

//OpenDB would open the db but the close has to be done after the action
func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(192.168.1.200:32769)/BikeTransport_db")
	if err != nil {
		fmt.Println("Failed to open SQL db:", err)
		log.Fatalln("Failed to open SQL db:", err)
	}
	fmt.Println("Database successfully opened.")

	return db
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
