package main

import (
	handler "Assignment1/Handler"
	"log"
	"net/http"
	"time"
)

func main() {
	port := "8080"
	path := "/librarystats/v1/"
	bookcountP := path + "bookcount/"
	readershipP := path + "readership/"
	statusP := path + "status/"

	http.HandleFunc(bookcountP, handler.Bookcount)
	http.HandleFunc(readershipP, handler.Readership)
	http.HandleFunc(statusP, handler.Status)

	StartTime := time.Now()
	handler.StartTime = StartTime

	log.Println("Listening on port " + port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

/*
	Tests:

	localhost:8080/librarystats/v1/bookcount/?language=no
   	localhost:8080/librarystats/v1/bookcount/?language=no,fi
	localhost:8080/librarystats/v1/bookcount/?language=en

   	localhost:8080/librarystats/v1/readership/no
	localhost:8080/librarystats/v1/readership/en/?limit=8

	localhost:8080/librarystats/v1/status
*/
