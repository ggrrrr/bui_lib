package tests

import (
	"log"

	"github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/scylladb/gocqlx/v2/qb"
)

// CREATE TABLE IF NOT EXISTS  person (
// 	namespace text,
// 	id uuid,
// 	pin text,
// 	name text,
// 	email text,
// 	first_name text,
// 	last_name text,
// 	other_names text,
// 	phones map<text,text>,
// 	attrs frozen<map<text,text>>,
// 	labels set<text>,
// 	created_time timestamp,
// 	PRIMARY KEY ((id), namespace )
// )

func iter(qb *qb.SelectBuilder, filterValues *qb.M) (*[]Person, error) {
	var out []Person
	q := qb.Query(*cassandra.Session).BindMap(*filterValues)
	if q.Err() != nil {
		log.Printf("stmt qb2.Err(): %v", q.Err())
		return nil, q.Err()

	}
	log.Printf("stmt qb: %v", q)
	iter := q.Iter()
	if iter == nil {
		log.Printf("ASDADSASDASDASDASDAS: %v", q.Err())
		return nil, q.Err()

	}
	var person Person
	for iter.StructScan(&person) {
		log.Printf("iter.Warnings() %+v:", iter.Warnings())

		log.Printf("%+v:", person)
		out = append(out, person)
	}
	if err := iter.Close(); err != nil {
		log.Printf("StructScan() failed:", err)
		return nil, err
	}
	return &out, nil
}
