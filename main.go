package main

import (
	"log"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

func init() {
	db, err := bbolt.Open("tickets.db", 0600, nil)
	act.Assert(err)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
