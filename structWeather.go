package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//JSONResult is exported
type JSONResult struct {
	Type     string    `json:"type"`
	Features []Content `json:"features"`
}

//Content is exported
type Content struct {
	Type     string  `json:"type"`
	Geometry Weather `json:"geometry"`
	Prop Rainfall  `json:"properties"`
}

//Weather is exported
type Weather struct {
	Type string    `json:"type"`
	Coor []float64 `json:"coordinates"`
}

//Rainfall is exported
type Rainfall struct {
	Rain float64 `json:"rain"`
}

//W is exported
type W struct {
	Lat      float64
	Long     float64
	Rainfall float64
}

func weatherHTML(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) { // check if user is already logged in
		//User is not signed in
		http.Redirect(res, req, "/", http.StatusSeeOther)
		Warning.Println("Unauthorised request.")
	}

	m := fmt.Sprintf("%d", time.Now().Month())
	d := fmt.Sprintf("%d", time.Now().Day())
	h := fmt.Sprintf("%d", time.Now().Hour())
	if len(m) == 1 {
		m = fmt.Sprintf("0%s", m)
	}
	if len(d) == 1 {
		d = fmt.Sprintf("0%s", d)
	}
	if len(h) == 1 {
		h = fmt.Sprintf("0%s", h)
	}
	t := fmt.Sprintf("%d%s%s%s0000", time.Now().Year(), m, d, h)
	jsonURL := fmt.Sprintf("https://nowcast.tk/rain/50km/" + t)
	var tempList []W
	resp, err := http.Get(jsonURL);
	if err == nil {
		defer resp.Body.Close()                             //After everything is done, close the connection, only close it if the connection is a success
		body, err := ioutil.ReadAll(resp.Body);
		if err == nil { //Returns the body obtained from the Web API and an err
			var result JSONResult
			json.Unmarshal(body, &result)
			if result.Type != "" { //Managed to fetch the result
				for _, item := range result.Features {
					w := W{item.Geometry.Coor[0], item.Geometry.Coor[1], item.Prop.Rain*10}
					tempList = append(tempList, w)
				}
			}
		}
	}

	tpl.ExecuteTemplate(res, "weather.gohtml", tempList)
}
