package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "stoyan:otednodoosem@tcp(127.0.0.1:3306)/bank?multiStatements=true")
	if err != nil {
		panic(err)
	}

	miroMoney := 0.0
	err = db.QueryRow("select money from users where name='miro'").Scan(&miroMoney)
	if err != nil {
		panic(err)
	}
	fmt.Printf("miro before: %f\n", miroMoney)

	var name = "'; update users set money=1000000000 where name='miro' -- "
	// var name = "miro"
	var amount float32 = 10
	query := fmt.Sprintf("update users set money=money-%f where name='%s'", amount, name)
	// query := "update users set money=money-? where name=?";
	fmt.Println("query: ", query)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	miroMoney = 0.0
	err = db.QueryRow("select money from users where name='miro'").Scan(&miroMoney)
	if err != nil {
		panic(err)
	}
	fmt.Printf("miro after: %f\n", miroMoney)
}
