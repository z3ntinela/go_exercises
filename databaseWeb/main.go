package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"time"
)

type person struct { // Con minúscula es una estructura privada, con mayúscula es publica, lo mismo pasa con funciones
	Id   int
	Name string
	Age  int
}

func dbConnect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", connStr)

	if err != nil {
		return db, err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return db, err
	}
	return db, nil
}

func getPerson(id int, db *sql.DB) (person, bool, error) {
	person := person{}
	sql := "select id, name, age from client where id=?"
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, sql, id)
	if err != nil {
		return person, false, err
	}

	if err != nil {
		return person, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var idPerson int
		var name string
		var age int
		err := rows.Scan(&idPerson, &name, &age)
		if err != nil {
			return person, false, err
		}
		person.Id = idPerson
		person.Name = name
		person.Age = age
		return person, true, nil
	}
	return person, false, nil
}

func main() {

	connStr := "./clients.db"

	db, err := dbConnect(connStr)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/person",
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				log.Fatal(err)
			}
			p, ok, err := getPerson(id, db)
			if err != nil {
				log.Fatal(err)
			}

			if ok {
				fmt.Println(p)
				jsonBytes, _ := json.Marshal(p)
				w.Write(jsonBytes)
			} else {

				w.Write([]byte("Persona no encontrada"))
			}
		},
	)

	webServer := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  90 * time.Second,
		Handler:      mux,
	}

	err = webServer.ListenAndServe()

	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
