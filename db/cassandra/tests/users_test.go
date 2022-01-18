package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/ggrrrr/bui_lib/db"
	"github.com/ggrrrr/bui_lib/db/cassandra"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

const (
	NS   string = "ns"
	PIN1 string = "PIN1"
)

var (
	PHONES1 = map[string]string{"label1": "val1"}
	LABEL1  = []string{"l1", "l2", "l1=1", "l1=2"}
)

func OK(t *testing.T, str string, v ...interface{}) {
	a := fmt.Sprintf(str, v...)
	t.Logf("OK WITH %v", a)
}

func TestPasswd1(t *testing.T) {
	var err error
	// ctx := context.Background()
	// t.Setenv(db.DB_CLUSTER, "127.0.0.1")
	// t.Setenv(db.DB_KEYSPACE, "test")
	ctx := context.WithValue(context.Background(), api.ContextPerson, NS)
	err = cassandra.Configure()
	if err != nil {

		t.Fatalf("CreateSession: %v", err)
	}
	_, err = cassandra.Connect()
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	err = db.CreateSchema("people")
	if err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	person1 := &Person{Namespace: NS}
	err = person1.Insert(ctx, cassandra.Session)
	if err == nil {
		t.Errorf("cant insert: %v", err)
	}
	t.Logf("OK ERROR:%v", err)

	person2 := &Person{Namespace: NS, Id: gocql.UUID{}}
	err = person2.Insert(ctx, cassandra.Session)
	if err == nil {
		t.Errorf("cant insert: %v", err)
	}
	OK(t, "err: %v", err)

	idNew, err := uuid.NewUUID()
	PIN1 := idNew.Domain().String()
	if err != nil {
		t.Fatalf("cant insert: %v, %v", err, idNew)
	}
	idNewBytes, err := idNew.MarshalBinary()
	t.Logf("OK UUID %v :%v", PIN1, idNew)
	if err != nil {
		t.Fatalf("cant insert: %v, %v", err, idNew)
	}
	OK(t, "OK UUID:%v", idNewBytes)
	uuid1, err := gocql.UUIDFromBytes(idNewBytes)

	if err != nil {
		t.Fatalf("cant insert: %v, %v", err, idNew)
	}

	person3 := &Person{Namespace: NS, Id: uuid1, PIN: PIN1, Labels: LABEL1, Phones: PHONES1}
	err = person3.Insert(ctx, cassandra.Session)
	if err != nil {
		t.Errorf("cant insert: %v", err)
	}

	res1, err := GetByPin(ctx, cassandra.Session, PIN1)

	if err != nil {
		t.Errorf("cant select by pin: %v", err)
	}
	OK(t, "p: %v", res1)

	res2, err := GetByUUID(ctx, cassandra.Session, uuid1)
	if err != nil {
		t.Fatalf("cant select by pin: %v", err)
	}
	if len(*res2) != 1 {
		t.Errorf(" too many results %v", res2)
	}
	OK(t, "GetByUUID: %v", res2)

	res3, err := GetByLabels(ctx, cassandra.Session, "asd")
	if err != nil {
		t.Fatalf("cant select by labels: %v", err)
	}
	OK(t, "GetByUUID: %v", res3)
}
