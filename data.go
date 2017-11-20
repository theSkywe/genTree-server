package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB struct
type DB struct {
	*sql.DB
}

const dbName = "Tree"
const dbCOnnectionString = "root:root@tcp(localhost:3306)/" + dbName

var scheme = `CREATE TABLE IF NOT EXISTS nodes (
				id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
				name varchar(255) NOT NULL,
				image varchar(255) NOT NULL,
				lft integer NOT NULL,
				rgt integer NOT NULL,
				depth integer NOT NULL);`

/* do we really need size of images?
type Image struct {
	Path string `json:"path"`
	Height int `json:"height"`
	Width int `json:"width"`
}
*/

// Node type struct
type Node struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Lft   int    `json:"lft,omitempty"`
	Rgt   int    `json:"rgt,omitempty"`
	Depth int    `json:"depth,omitempty"`
}

//Nodes type
type Nodes []Node

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewOpen opens connections with database
func newOpen(dt string, conn string) (DB, error) {
	db, err := sql.Open(dt, conn)
	return DB{db}, err
}

// InitDB create DB and table if not exists
func initDB(name string) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
	checkErr(err)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	checkErr(err)

	_, err = db.Exec("USE " + dbName)
	checkErr(err)

	_, err = db.Exec(scheme)
	checkErr(err)

	initRoot := "INSERT INTO nodes(name, image, lft, rgt, depth) VALUES(?, ?, 1, 2, 1); "
	initName := "You"
	initImage := "../images/you.png"
	_, err = db.Exec(initRoot, initName, initImage)
	checkErr(err)
}

// GetAllNodesFromDB is a function that recieving all nodes from database into struct Nodes
func (db DB) getAllNodesFromDB() (nodes Nodes, err error) {
	query := `SELECT node.id, node.name, node.image, node.lft, node.rgt, node.depth
				FROM nodes as node,
					nodes as parent
				WHERE node.lft BETWEEN parent.lft and parent.rgt
					AND parent.id = node.id
				GROUP BY node.id
				ORDER BY node.lft`
	rows, err := db.Query(query)
	checkErr(err)
	var ns Nodes
	nodes = ns
	for rows.Next() {
		var n Node
		err := rows.Scan(&n.ID, &n.Name, &n.Image, &n.Lft, &n.Rgt, &n.Depth)
		checkErr(err)
		nodes = append(nodes, n)
	}
	rows.Close()
	return
}

// InsertNodeIntoDB is a function that add new node from database
func (db DB) insertNodeIntoDB(parentID int, name string, image string) (err error) {
	tx, err := db.Begin()
	checkErr(err)

	var lockTable = "LOCK TABLE nodes WRITE"
	var selectVal = `	SELECT nodes.lft, nodes.rgt, nodes.depth
						FROM nodes
						WHERE id = ?`
	var updateRight = "UPDATE nodes SET rgt = rgt + 2 WHERE rgt > ?"
	var updateLeft = "UPDATE nodes SET lft = lft + 2 WHERE lft > ?"
	var insert = "INSERT INTO nodes(name, image, lft, rgt, depth) VALUES(?, ?, ? + 1, ? + 2, ? + 1)"
	var unlockTable = "UNLOCK TABLES"

	var lft, rgt, dep, val int

	tx.Exec(lockTable)
	err = tx.QueryRow(selectVal, parentID).Scan(&lft, &rgt, &dep)
	checkErr(err)

	if (rgt - lft) < 1 {
		val = rgt
	} else {
		val = lft
	}

	tx.Exec(updateRight, val)
	tx.Exec(updateLeft, val)

	tx.Exec(insert, name, image, val, val, dep)

	tx.Exec(unlockTable)
	tx.Commit()

	return
}

// DeleteNodeFromDB is a function that delete node from database
func (db DB) deleteNodeFromDB(id int) (err error) {
	tx, err := db.Begin()
	checkErr(err)

	var lockTable = "LOCK TABLE nodes WRITE"
	var selectVal = `	SELECT nodes.lft, nodes.rgt
						FROM nodes
						WHERE id = ?`
	var delete = "DELETE FROM nodes WHERE lft BETWEEN ? AND ?"
	var updateRight = "UPDATE nodes SET rgt = rgt - ? WHERE rgt > ?"
	var updateLeft = "UPDATE nodes SET lft = lft - ? WHERE lft > ?"
	var unlockTable = "UNLOCK TABLES"

	var lft, rgt, width int

	tx.Exec(lockTable)
	err = tx.QueryRow(selectVal, id).Scan(&lft, &rgt)
	checkErr(err)

	width = rgt - lft + 1

	tx.Exec(delete, lft, rgt)

	tx.Exec(updateRight, width, rgt)
	tx.Exec(updateLeft, width, rgt)

	tx.Exec(unlockTable)
	tx.Commit()

	return
}
