package outletapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"outletapi/outletconfiguration"
	"outletapi/outletsecurity"
	"outletapi/outletwebserver"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var secureCookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var securePages = make(map[string]bool)

var config outletconfiguration.Configuration

const websiteName = "Outlet API - "

// SetupEndpoints defines API routes
func SetupEndpoints(router *mux.Router) {
	config = outletconfiguration.LoadConfiguration()

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/status", apiStatus).Methods("GET")
	router.HandleFunc("/login", loginUser).Methods("POST")
	router.HandleFunc("/logout", logoutUser).Methods("POST")
	router.HandleFunc("/turn", turn).Methods("GET")
	router.HandleFunc("/turn", turnSwitch).Methods("POST")

	defineSecurePages()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), router))
}

func defineSecurePages() {
	securePages["turn"] = true
}

func homePage(w http.ResponseWriter, req *http.Request) {
	page := outletwebserver.LoadPage("Index")

	fmt.Fprintf(w, "<title> %s </title> <div> %s </div>", websiteName+"Home", page.Content)
}

func apiStatus(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("I'm doing okay")
}

func loginUser(w http.ResponseWriter, req *http.Request) {
	var u outletsecurity.User
	json.NewDecoder(req.Body).Decode(&u)

	if outletsecurity.CheckUser(u) {
		setCookie(u.Username, w)
	} else {
		removeCookie(w)
	}
}

func logoutUser(w http.ResponseWriter, req *http.Request) {
	removeCookie(w)
	http.Redirect(w, req, "/", 302)
}

func turn(w http.ResponseWriter, req *http.Request) {
	authorize(w, req, "turn")

	page := outletwebserver.LoadPage("turn")
	fmt.Fprintf(w, "<title> %s </title> <div> %s </div>", websiteName+"Turn", page.Content)
}

func turnSwitch(w http.ResponseWriter, req *http.Request) {
	authorize(w, req, "turn")

	switchState := "unknown"
	value := req.FormValue("value")

	if value != "" {
		jsonValue, error := json.Marshal(map[string]string{
			"value1": value,
		})
		if error == nil {
			response, _ := http.Post(config.VeSyncEndpoint,
				"application/json", bytes.NewBuffer(jsonValue))
			if response.Status == "200 OK" {
				switchState = value
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"State": switchState,
	})
}

func setCookie(username string, response http.ResponseWriter) {
	value := map[string]string{
		"name": username,
	}
	encoded, error := secureCookieHandler.Encode("outlet", value)
	if error == nil {
		cookie := &http.Cookie{Name: "outlet", Value: encoded}

		http.SetCookie(response, cookie)
	}
}

func removeCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{Name: "outlet", Value: ""}

	http.SetCookie(w, cookie)
}

func readCookie(request *http.Request) (username string) {
	cookie, error := request.Cookie("outlet")
	if error == nil {
		cookieValue := make(map[string]string)
		error := secureCookieHandler.Decode("outlet", cookie.Value, &cookieValue)
		if error == nil {
			username = cookieValue["name"]
		}
	}

	return username
}

func authorize(response http.ResponseWriter, request *http.Request, page string) {
	if securePages[page] {
		cookieValue := readCookie(request)
		if cookieValue != "" && len(cookieValue) > 3 {
			return
		}
	}
	http.Redirect(response, request, "/", 302)
}
