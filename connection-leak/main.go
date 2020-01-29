package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "stoyan:otednodoosem@tcp(localhost:3306)/bank")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(2)
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			rows, err := db.Query("select name from users")
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			for rows.Next() {
				var name string
				err := rows.Scan(&name)
				if err != nil {
					panic(err)
				}
				fmt.Printf("name: %s\n", name)
				time.Sleep(2 * time.Second)
			}
			if err = rows.Err(); err != nil {
				panic(err)
			}
			fmt.Println("all user read")
		}()
	}
	wg.Wait()
}
