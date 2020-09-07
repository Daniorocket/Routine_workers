package numbersql

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDb(tablename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./"+tablename+".db")
	if err != nil {
		return nil, err
	}
	return db, err
}
func InitDb(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS numbers (numberid INTEGER PRIMARY KEY AUTOINCREMENT,workername varchar(64) NULL,number INTEGER)")
	stmt.Exec()
	fmt.Println("Init db done!")
	return err
}
func InsertRow(db *sql.DB, workername string, number int) error {
	stmt, err := db.Prepare("INSERT INTO numbers(workername, number) values(?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(workername, number)
	if err != nil {
		return err
	}
	res.LastInsertId()
	return nil
}
func SelectAllData(db *sql.DB) error {
	rows, err := db.Query("SELECT numberid,workername,number FROM numbers")
	if err != nil {
		return err
	}
	var numberid int
	var workername string
	var number string
	fmt.Println("SELECT:")
	defer rows.Close() //good habit to close
	for rows.Next() {
		if err = rows.Scan(&numberid, &workername, &number); err != nil {
			return err
		}
		fmt.Println("Record:", numberid, workername, number)
	}
	return nil
}
