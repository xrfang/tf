package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

func save(db *bbolt.DB) {
	t := Ticket{
		Entries: []*Content{{
			Caption: "This is a test",
			Type:    "text/plain",
			Data:    []byte("Hello, world"),
			Creator: 1,
			Created: time.Now(),
			Updated: time.Now(),
		}},
		Due:      time.Now().Add(48 * time.Hour).UnixNano(),
		Assignee: 1,
	}
	act.Assert(t.Save(db))
	je := json.NewEncoder(os.Stdout)
	je.SetIndent("> ", "    ")
	je.Encode(t)
}

func load(db *bbolt.DB) {
	var t Ticket
	act.Assert(t.Load(db, 1))
	je := json.NewEncoder(os.Stdout)
	je.SetIndent("< ", "    ")
	je.Encode(t)
}

func main() {
	db, err := bbolt.Open("tickets.db", 0600, nil)
	act.Assert(err)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	save(db)
	load(db)
	fmt.Printf("done\n")
}
