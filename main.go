package main

import (
	"fmt"
	"log"
	"time"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

func main() {
	db, err := bbolt.Open("tickets.db", 0600, nil)
	act.Assert(err)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bktTickets)
		b.ForEach(func(k, v []byte) error {
			id, err := decode[uint64](k)
			act.Assert(err)
			fmt.Printf("key=%v; val=%s\n", id, string(v))
			return nil
		})
		return nil
	})
	fmt.Printf("done\n")
	return
	t := Ticket{
		Entries: []Content{{
			Caption: "This is a test",
			Type:    "text/plain",
			Data:    []byte("Hello, world"),
			Creator: 1,
			Created: time.Now(),
			Updated: time.Now(),
		}},
		Due:      time.Now().Add(48 * time.Hour),
		Assignee: 1,
	}
	act.Assert(t.Save(db))
}
