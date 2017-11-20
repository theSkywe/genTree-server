package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getTreeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := newOpen("mysql", dbCOnnectionString)
	checkErr(err)
	defer db.Close()

	nodes, err := db.getAllNodesFromDB()
	checkErr(err)
	json.NewEncoder(w).Encode(nodes)
}

func addNodeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := newOpen("mysql", dbCOnnectionString)
	checkErr(err)
	defer db.Close()

	var n Node
	_ = json.NewDecoder(r.Body).Decode(&n)
	checkErr(err)

	err = db.insertNodeIntoDB(n.ID, n.Name, n.Image)
	checkErr(err)

	nodes, err = db.getAllNodesFromDB()
	checkErr(err)

	json.NewEncoder(w).Encode(nodes)
}

func deleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db, err := newOpen("mysql", dbCOnnectionString)
	checkErr(err)
	defer db.Close()

	nodeID, err := strconv.Atoi(params["id"])
	checkErr(err)

	err = db.deleteNodeFromDB(nodeID)
	checkErr(err)

	nodes, err = db.getAllNodesFromDB()
	checkErr(err)

	json.NewEncoder(w).Encode(nodes)
}
