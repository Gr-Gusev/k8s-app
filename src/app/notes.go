package main

import (
	"database/sql"
	"log"
	"net/http"
)

type Note struct {
	Id   int64
	Text string
}

func GetAll(db *sql.DB) ([]Note, error) {
	rows, err := db.Query("select * from notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		note := Note{}
		err := rows.Scan(&note.Id, &note.Text)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func Add(db *sql.DB, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	text := r.PostForm.Get("text")
	if text != "" {
		res, err := db.Exec("insert into notes (text) values (?)", text)
		if err != nil {
			return err
		}

		id, _ := res.LastInsertId()
		log.Printf("successfully insert new value, last inserted id: %d", id)
	}

	return nil
}

func Delete(db *sql.DB, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	id := r.PostForm.Get("id")
	if id != "" {
		_, err := db.Exec("delete from notes where id=?", id)
		if err != nil {
			return err
		}

		log.Printf("successfully delete note with id: %v", id)
	}

	return nil
}
