package main

import (
	"fmt"
	"log"
	"my_api/api"
	"net/http"
)

func main() {
	//get a new router
	r := api.NewRouter()
	fmt.Println("Running...")
	//serve on port 1701
	log.Fatal(http.ListenAndServe(":1701", r))
}
