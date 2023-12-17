package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
)

type Pdb struct {
	dbObj *sql.DB
	table string
}

func (db *Pdb) DbObj() *sql.DB {
	return db.dbObj
}

func (db *Pdb) SetDbObj(dbObj *sql.DB) error {
	db.dbObj = dbObj
	return nil
}

func getDb(ymlFile string) *sql.DB {
	data := getEnv(ymlFile)
	return dbConn(data)
}

func getEnv(yName string) map[interface{}]interface{} {
	ymlFile, err := ioutil.ReadFile(yName)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})
	err2 := yaml.Unmarshal(ymlFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}
	return data
}

func dbConn(data map[interface{}]interface{}) *sql.DB {
	var connStr string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		data["host"], data["port"], data["user"], data["password"], data["database"])

	dbObj, err := sql.Open("postgres", connStr)
	checkError(err)
	return dbObj
}

func (db *Pdb) Table() string {
	return db.table
}

func (db *Pdb) SetTable(table string) error {
	db.table = table
	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func Setdb() {
	pgDb := Pdb{}
	pgDb.SetDbObj(getDb("./db/db.yaml"))
	pgDb.SetTable("id")

	defer pgDb.dbObj.Close()
	dbPing(pgDb.dbObj)

	//insertTbl(pgDb)

	rows := selectTbl(pgDb)
	defer rows.Close()

	//printRows(rows)
}

func printRows(rows *sql.Rows) {
	var id string
	var server string
	var computer string

	switch err := rows.Scan(&id, &server, &computer); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
	case nil:
		fmt.Printf("%s, %s, %s", id, server, computer)
	default:
		checkError(err)
	}
}

func selectTbl(pgDb Pdb) *sql.Rows {
	stmt1 := fmt.Sprintf("SELECT * from %s;", pgDb.table)
	rows, err := pgDb.dbObj.Query(stmt1)
	checkError(err)
	return rows
}

func insertTbl(pgDb Pdb) {
	stmtIns := fmt.Sprintf("INSERT INTO %s (id, computer) VALUES ($1, $2);", pgDb.table)
	_, err := pgDb.dbObj.Exec(stmtIns, "sukho8757@gmail.com", 58)
	checkError(err)
	_, err = pgDb.dbObj.Exec(stmtIns, "joke@nana.com", 57)
	checkError(err)
}

func createTbl(pgDb Pdb) {
	creStr := fmt.Sprintf("CREATE TABLE %s (id serial PRIMARY KEY, name VARCHAR(20), quantity INTEGER);", pgDb.table)
	_, err := pgDb.dbObj.Exec(creStr)
	checkError(err)
	fmt.Println("Finished creating table")
}

func dropTable(pgDb Pdb) {
	dropStr := fmt.Sprintf("DROP TABLE IF EXISTS %s;", pgDb.table)
	_, err := pgDb.dbObj.Exec(dropStr)
	checkError(err)
	fmt.Println("Finished dropping table (if existed)")
}

func dbPing(db *sql.DB) {
	err := db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")
}

func AddLog(jsonData map[string]interface{}) {
	pgDb := Pdb{}
	pgDb.SetDbObj(getDb("./db/db.yaml"))
	defer pgDb.dbObj.Close()
	stmtIns := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 AND server=$2;", "id")
	rows, err := pgDb.dbObj.Query(stmtIns, jsonData["id"], jsonData["server"])
	checkError(err)
	if rows.Next() {
		var (
			index    int
			id       string
			server   string
			done     bool
			date     time.Time
			computer string
		)
		rows.Scan(&index, &id, &server, &done, &date, &computer)
		if computer != jsonData["computer"] {
			stmtIns2 := fmt.Sprintf("UPDATE %s SET computer=$3 WHERE id=$1 AND server=$2;", "id")
			_, err := pgDb.dbObj.Exec(stmtIns2, jsonData["id"], jsonData["server"], jsonData["computer"])
			checkError(err)
		}
		stmtIns3 := fmt.Sprintf("INSERT INTO %s (id_index, code, detail, date) VALUES ($1, $2, $3, $4);", "log")
		_, err := pgDb.dbObj.Exec(stmtIns3, index, jsonData["code"], jsonData["detail"], time.Now())
		checkError(err)
	} else {
		stmtIns2 := fmt.Sprintf("INSERT INTO %s (id, server, done, date, computer) VALUES ($1, $2, $3, $4, $5);", "id")
		_, err := pgDb.dbObj.Exec(stmtIns2, jsonData["id"], jsonData["server"], false, time.Now(), jsonData["computer"])
		checkError(err)
	}
	defer rows.Close()

}
