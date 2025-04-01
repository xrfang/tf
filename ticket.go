package main

import (
	"os"

	"go.etcd.io/bbolt"
	"go.xrfang.cn/act"
)

type (
	Ticket struct {
		id       uint64
		Entries  []*Content
		Status   int8
		Tags     []string
		Metrics  map[string]float64
		Due      int64 //time.UnixNano
		Assignee uint64
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
	tb.DeleteBucket(bktContent)
	var cb *bbolt.Bucket
	act.Assure(tb.CreateBucket(bktContent)).Scan(&cb)
	for i, e := range t.Entries {
		act.Assert(cb.Put(encode(uint16(i)), e.encode()))
	}
}

func (t *Ticket) saveMetrics(tb *bbolt.Bucket) {
	tb.DeleteBucket(bktMetrics)
	if len(t.Metrics) > 0 {
		var mb *bbolt.Bucket
		act.Assure(tb.CreateBucket(bktMetrics)).Scan(&mb)
		for k, v := range t.Metrics {
			bput(mb, k, v)
		}
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
	bput(tb, "due", t.Due)
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

func (t *Ticket) loadContents(tb *bbolt.Bucket) {
	if cb := tb.Bucket(bktContent); cb != nil {
		cb.ForEach(func(_, v []byte) error {
			var e Content
			t.Entries = append(t.Entries, e.decode(v))
			return nil
		})
	} else {
		t.Entries = []*Content{}
	}
}

func (t *Ticket) loadMetrics(tb *bbolt.Bucket) {
	if mb := tb.Bucket(bktMetrics); mb != nil {
		mb.ForEach(func(k, v []byte) error {
			key := act.Assure(decode[string](k))[0].([]string)[0]
			val := act.Assure(decode[float64](v))[0].([]float64)[0]
			t.Metrics[key] = val
			return nil
		})
	} else {
		t.Metrics = map[string]float64{}
	}
}

func (t *Ticket) load(b *bbolt.Bucket, id uint64) (err error) {
	tb := b.Bucket(encode(id))
	if tb == nil {
		return os.ErrNotExist
	}
	defer act.Catch(&err)
	t.loadContents(tb)
	t.loadMetrics(tb)
	t.Status = 0
	if s := bget[int8](tb, "status"); s != nil {
		t.Status = s[0]
	}
	t.Tags = bget[string](tb, "tags")
	t.Due = 0
	if d := bget[int64](tb, "due"); len(d) > 0 {
		t.Due = d[0]
	}
	t.Assignee = 0
	if a := bget[uint64](tb, "assignee"); len(a) > 0 {
		t.Assignee = a[0]
	}
	t.Items = bget[uint64](tb, "items")
	t.id = id
	return nil
}

func (t *Ticket) Load(db *bbolt.DB, id uint64) (err error) {
	defer act.Catch(&err)
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bktTickets)
		if b == nil {
			return os.ErrNotExist
		}
		return t.load(b, id)
	})
}
