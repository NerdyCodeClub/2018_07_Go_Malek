package main

import (
	"fmt"
	"outletapi/api"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var router = mux.NewRouter()
var secureCookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func main() {
	fmt.Println("Starting the API.")
	outletapi.SetupEndpoints(router)
}
