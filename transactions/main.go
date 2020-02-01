package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "stoyan:otednodoosem@tcp(127.0.0.1:3306)/bank")
	if err != nil {
		panic(err)
	}

	err = donate(db, "stoyan", "miro", 10.0)
	if err != nil {
		fmt.Printf("failed to donate: %s\n", err)
	}

	miroMoney := 0.0
	err = db.QueryRow("select money from users where name='miro'").Scan(&miroMoney)
	if err != nil {
		panic(err)
	}

	stoyanMoney := 0.0
	err = db.QueryRow("select money from users where name='stoyan'").Scan(&stoyanMoney)
	if err != nil {
		panic(err)
	}
	fmt.Printf("stoyan: %f, miro: %f", stoyanMoney, miroMoney)
}

func donate(db *sql.DB, from string, to string, amount float32) error {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("update users set money=money - ? where name=?", amount, from)
	if err != nil {
		return err
	}

	_, err = tx.Exec("update users set money=money + ? where name=?", amount, to)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
