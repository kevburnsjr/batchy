package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kevburnsjr/batchy"
)

var db *sql.DB

func main() {
	/*
		DROP DATABASE IF EXISTS testdb;
		CREATE DATABASE testdb;
		GRANT ALL PRIVILEGES ON testdb.* TO 'testuser'@'localhost' IDENTIFIED BY 'testpass';
		CREATE TABLE `testdb`.`test` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, `uid` VARCHAR(255) NULL, PRIMARY KEY (`id`));
	*/
	var err error
	db, err = sql.Open("mysql", "testuser:testpass@/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/unbatched", func(w http.ResponseWriter, r *http.Request) {
		err := write(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	http.HandleFunc("/batched", func(w http.ResponseWriter, r *http.Request) {
		err := writeBatched(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func write(str string) (err error) {
	_, err = db.Exec(`INSERT INTO test (uid) VALUES (?)`, str)
	return
}

func writeBatched(str string) (err error) {
	err = batcher.Add(str)
	return
}

var batcher = batchy.New(50, 100*time.Millisecond, func(items []interface{}) (errs []error) {
	q := fmt.Sprintf(`INSERT INTO test (uid) VALUES %s`, strings.Trim(strings.Repeat(`(?),`, len(items)), ","))
	_, err := db.Exec(q, items...)
	if err != nil {
		errs = make([]error, len(items))
		for i := range errs {
			errs[i] = err
		}
	}
	return
})
