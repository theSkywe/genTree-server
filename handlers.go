package main

import (
	"encoding/json"
	"net/http"
)

func getTreeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := newOpen("mysql", dbCOnnectionString)
	checkErr(err)
	defer db.Close()

	nodes, err := db.getAllNodesFromDB()
	checkErr(err)
	json.NewEncoder(w).Encode(nodes)
}

func addNodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//params := mux.Vars(r)
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
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	db, err := newOpen("mysql", dbCOnnectionString)
	checkErr(err)
	defer db.Close()

	var n Node
	_ = json.NewDecoder(r.Body).Decode(&n)

	err = db.deleteNodeFromDB(n.ID)
	checkErr(err)

	nodes, err = db.getAllNodesFromDB()
	checkErr(err)

	json.NewEncoder(w).Encode(nodes)
}
