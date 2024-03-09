package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("New account number: ", acc.Number)
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "John", "Smith", "johnsmith")
}

func main() {

	//Seed
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	if *seed {
		fmt.Println("Seeding the database")
	}
	seedAccounts(store)

	//fmt.Printf("%+v\n", store)
	server := NewAPIServer(":3000", store)
	server.Run()
}
