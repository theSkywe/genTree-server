package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var nodes Nodes

func main() {
	db, err := newOpen("mysql", "root:root@tcp(localhost:3306)/")
	checkErr(err)
	defer db.Close()

	initDB(dbName)

	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"Content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"})

	router.HandleFunc("/nodes", getTreeHandler).Methods("GET")
	router.HandleFunc("/nodes", addNodeHandler).Methods("POST")
	router.HandleFunc("/nodes/{id}", deleteNodeHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
