package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const uploadedPath = "./uploaded/"

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

	r.ParseMultipartForm(1000000)
	form := r.MultipartForm

	var n Node

	n.ID, _ = strconv.Atoi(r.FormValue("id"))
	n.Name = r.FormValue("name")

	files := form.File["image"]
	imageFile, err := files[0].Open()
	defer imageFile.Close()

	dst, err := os.Create(uploadedPath + files[0].Filename)
	checkErr(err)
	defer dst.Close()

	io.Copy(dst, imageFile)

	imagePath := uploadedPath + files[0].Filename
	checkErr(err)

	err = db.insertNodeIntoDB(n.ID, n.Name, imagePath)
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
