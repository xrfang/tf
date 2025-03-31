package main

import (
	"time"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

type (
	Ticket struct {
		id       uint64
		Entries  []Content
		Status   int8
		Tags     []string
		Metrics  map[string]float64
		Due      time.Time
		Assignee int64
		Items    []uint64
	}
)

func (t *Ticket) next(tx *bbolt.Tx) {
	rb := tx.Bucket(bktTickets)
	act.Assure(rb.NextSequence()).Scan(&t.id)
}

func (t *Ticket) ID() uint64 {
	return t.id
}

func (t *Ticket) saveContents(tb *bbolt.Bucket) {
	var cb *bbolt.Bucket
	act.Assure(tb.CreateBucket(bktContent)).Scan(&cb)
	for _, e := range t.Entries {
		bput(cb)
	}
}

func (t *Ticket) saveMetrics(tb *bbolt.Bucket) {
	var mb *bbolt.Bucket
	act.Assure(tb.CreateBucket(bktMetrics)).Scan(&mb)
	for k, v := range t.Metrics {
		bput(mb, k, v)
	}
}

func (t *Ticket) save(b *bbolt.Bucket) {
	if t.id == 0 {
		t.next(b.Tx())
	}
	var tb *bbolt.Bucket
	act.Assure(b.CreateBucketIfNotExists(encode(t.ID()))).Scan(&tb)
	bput(tb, "id", t.ID())
	bput(tb, "status", t.Status)
	bput(tb, "tags", t.Tags...)
	if len(t.Entries) > 0 {
		t.saveContents(tb)
	}
	if len(t.Metrics) > 0 {
		t.saveMetrics(tb)
	}
	bput(tb, "due", t.Due.UnixNano())
	bput(tb, "assignee", t.Assignee)
	bput(tb, "items", t.Items...)
}

func (t *Ticket) Save(db *bbolt.DB) (err error) {
	defer act.Catch(&err)
	act.Assert(db.Update(func(tx *bbolt.Tx) error {
		b := act.Assure(tx.CreateBucketIfNotExists(bktTickets))[0].(*bbolt.Bucket)
		t.save(b)
		return nil
	}))
	return nil
}
