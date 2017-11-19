package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var nodes Nodes

func main() {
	db, err := newOpen("mysql", "root:root@tcp(localhost:3306)/")
	checkErr(err)
	defer db.Close()

	initDB(dbName)

	router := mux.NewRouter()
	router.HandleFunc("/nodes", getTreeHandler).Methods("GET")
	router.HandleFunc("/nodes", addNodeHandler).Methods("POST")
	router.HandleFunc("/nodes", deleteNodeHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
