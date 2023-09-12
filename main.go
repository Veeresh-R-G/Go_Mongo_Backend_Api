package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Veeresh-R-G/mongoapi/router"
)

func main() {
	fmt.Println("Mongo API")

	fmt.Println("Starting Application ... ")
	time.Sleep(2000)
	r := router.Router()

	log.Fatal(http.ListenAndServe(":5000", r))

	fmt.Println("Application Running on port 5000 ")

}
