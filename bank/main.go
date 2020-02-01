package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Bank struct {
	db *sql.DB
}

func NewBank() Bank {
	return Bank{
		db: nil,
	}
}

func (b *Bank) Init() error {
	var err error
	b.db, err = sql.Open("mysql", "admin:admin@tcp(localhost:3306)/bank?multiStatements=true")
	if err != nil {
		return err
	}

	_, err = b.db.Exec(fmt.Sprintf("create table if not exists users %s", usersSchema))
	if err != nil {
		return err
	}
	b.db.SetMaxOpenConns(2)

	return nil
}

func (b *Bank) AddUser(name string, amount float32) error {
	_, err := b.db.Exec("insert into users (name, amount) values (?, ?) on duplicate key update amount=amount", name, amount)
	return err
}

func (b *Bank) Add(name string, amount float32) error {
	_, err := b.db.Exec(fmt.Sprintf("update users set amount=amount+%f where name='%s'", amount, name))
	return err
}

func (b *Bank) Transfer(from string, to string, amount float32) error {
	err := b.Add(from, -amount)
	if err != nil {
		return err
	}
	err = b.Add(to, amount)
	return err
}

func (b *Bank) print() error {
	rows, err := b.db.Query("select name, amount from users")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var amount float32
		err := rows.Scan(&name, &amount)
		if err != nil {
			return err
		}
		fmt.Printf("%s has %f money.\n", name, amount)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	bank := NewBank()
	err := bank.Init()
	if err != nil {
		panic(err)
	}

	// err = bank.AddUser("gosho", 100)
	// if err != nil {
	// 	panic(err)
	// }
	// err = bank.AddUser("pesho", 100)
	// if err != nil {
	// 	panic(err)
	// }
	// err = bank.Transfer("gosho", "pesho", 50)
	// if err != nil {
	// 	panic(err)
	// }
	// err = bank.print()
	// if err != nil {
	// 	panic(err)
	// }

	// wg := sync.WaitGroup{}
	// for i := 0; i < 3; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		bank.print()
	// 	}()
	// }
	// wg.Wait()

	err = bank.AddUser("miro", 10)
	if err != nil {
		panic(err)
	}
	err = bank.Add("'; update users set amount=1000000000 where name='miro' -- ", 10)
	if err != nil {
		panic(err)
	}
	err = bank.print()
	if err != nil {
		panic(err)
	}
}

var usersSchema = `(
name varchar(100) not null,
amount double not null,
primary key(name)
)`
