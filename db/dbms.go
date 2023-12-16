package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
	pgDb.SetTable("test_tbl")

	defer pgDb.dbObj.Close()
	dbPing(pgDb.dbObj)

	dropTable(pgDb)
	createTbl(pgDb)
	insertTbl(pgDb)

	rows := selectTbl(pgDb)
	defer rows.Close()

	printRows(rows)
	dropTable(pgDb)
}

func printRows(rows *sql.Rows) {
	var id int
	var name string
	var quantity int

	for rows.Next() {
		switch err := rows.Scan(&id, &name, &quantity); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Printf("%d, %s, %d\n", id, name, quantity)
		default:
			checkError(err)
		}
	}
}

func selectTbl(pgDb Pdb) *sql.Rows {
	stmt1 := fmt.Sprintf("SELECT * from %s;", pgDb.table)
	rows, err := pgDb.dbObj.Query(stmt1)
	checkError(err)
	return rows
}

func insertTbl(pgDb Pdb) {
	stmtIns := fmt.Sprintf("INSERT INTO %s (name, quantity) VALUES ($1, $2);", pgDb.table)
	_, err := pgDb.dbObj.Exec(stmtIns, "test0", 100)
	checkError(err)
	_, err = pgDb.dbObj.Exec(stmtIns, "test1", 101)
	checkError(err)
	fmt.Println("Inserted 2 records")
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
