package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "stoyan:otednodoosem@tcp(localhost:3306)/shop")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	res, err := tx.Exec("insert into users (name, sex, country) values ('maria', 'F', 'Bulgaria')")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	userID, _ := res.LastInsertId()

	res, err = tx.Exec("insert into products (type, price) values ('racket', 300)")
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	productID, _ := res.LastInsertId()

	res, err = tx.Exec("insert into orders (user_id, product_id) values (?,?)", userID, productID)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("delete from users where id=?", userID)
	if err != nil {
		fmt.Println("failed to delete")
		tx.Rollback()
		// panic(err)

	}
	tx.Commit()
	fmt.Printf("inserted user: %d, product: %d\n", userID, productID)
	fmt.Println("Done")
}
