package main

import (
	"strings"
	"time"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

type Ticket struct {
	id       uint64
	Entries  []Content //sub-bucket
	Status   int
	Tags     []string
	Metrics  map[string]float64 //json
	Due      time.Time
	Assignee int64    //user-id
	Items    []Ticket //sub-bucket
}

func (t *Ticket) ID() uint64 {
	return t.id
}

func (t *Ticket) next(tx *bbolt.Tx) uint64 {
	var b *bbolt.Bucket
	act.Assure(tx.CreateBucketIfNotExists([]byte(bktTickets))).Scan(&b)
	return act.Assure(b.NextSequence())[0].(uint64)
}

func (t *Ticket) save(tx *bbolt.Tx) {
	if t.id == 0 {
		t.id = t.next(tx)
	}
	b := tx.Bucket([]byte(bktTickets))
	act.Assure(b.CreateBucket(encode(t.ID()))).Scan(&b)
	bput(b, "id", t.ID())
	bput(b, "status", t.Status)
	bput(b, "tags", strings.Join(t.Tags, "\n"))
	// 	Metrics  map[string]float64 //json
	// 	Due      time.Time
	// 	Assignee int64    //user-id
	// 	Items    []Ticket //sub-bucket
	// }
}

func (t *Ticket) Save(db *bbolt.DB) (err error) {
	defer act.Catch(&err)
	act.Assert(db.Update(func(tx *bbolt.Tx) error {
		t.save(tx)
		return nil
	}))
	return nil
}
