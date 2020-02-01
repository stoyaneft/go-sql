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

func (b *Bank) Deposit(name string, amount float32) error {
	_, err := b.db.Exec("update users set amount=amount+? where name=?", amount, name)
	return err
}

func (b *Bank) Transfer(from string, to string, amount float32) error {
	err := b.Deposit(from, -amount)
	if err != nil {
		return err
	}
	return b.Deposit(to, amount)
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
	return rows.Err()
}

func main() {
	bank := NewBank()
	err := bank.Init()
	if err != nil {
		panic(err)
	}
	err = bank.AddUser("miro", 10)
	if err != nil {
		panic(err)
	}

	err = bank.Deposit("'; update users set amount=10000000000 where name='miro' -- ", 10)
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
