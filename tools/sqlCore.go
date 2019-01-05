package tools

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"qiniupkg.com/x/log.v7"
)


var db_driver *sql.DB

func InitSqlDB(){
	os.Remove("./weedfs.db")
	db_driver, err := sql.Open("sqlite3", "./weedfs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db_driver.Close()

	sqlStmt := `
	create table info (id integer not null primary key, fid text, FileName text,created text,FileSize text);
	delete from info;
	`
	_, err = db_driver.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}


func SaveInfo(ID string, created string, FileName string, FileSize int32) int64 {
	db, err := sql.Open("sqlite3", "./weedfs.db")
	checkErr(err)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	//插入数据
	stmt, err := tx.Prepare("INSERT INTO info(fid, created ,FileName,FileSize ) values(?,?,?,?)")
	checkErr(err)

	defer stmt.Close()

	res, err := stmt.Exec(ID, created, FileName, FileSize)
	checkErr(err)
	tx.Commit()

	id, err := res.LastInsertId()
	checkErr(err)
	defer  db.Close()

	return id;
}


type TableInfo struct {
	Fid string
	Created string
	FileName string
	FileSize string
}

func QueryList()[]TableInfo {
	db, err := sql.Open("sqlite3", "./weedfs.db")
	checkErr(err)
	rows, err := db.Query("select fid, created ,FileName,FileSize from info")
	if err != nil {
		log.Fatal(err)
	}
	var TabList = []TableInfo{}

	defer rows.Close()
	for rows.Next() {
		var fid string
		var created string
		var FileName string
		var FileSize string
		err = rows.Scan(&fid, &created, &FileName, &FileSize)
		if err != nil {
			log.Fatal(err)
		}
		TabList = append(TabList,TableInfo{fid, created, FileName, FileSize})
	}
	return TabList;
}
func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
