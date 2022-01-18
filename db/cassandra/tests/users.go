package tests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ggrrrr/bui_lib/api"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

/*
CREATE TABLE person (
	namespace text,
	id uuid,
	name text,
	email text,
	first_name text,
	last_name text,
	phones frozen<map<text,text>>,
	attr frozen<map<text,text>>,
	labels frozen<map<text,text>>,
	created_time timestamp,
	PRIMARY KEY (namespace, id)
)
*/

type Person struct {
	Namespace   string `db:"namespace"`
	Id          gocql.UUID
	PIN         string
	Email       string
	Name        string
	FirstName   string
	LastName    string
	OtherNames  string
	Phones      map[string]string
	Labels      []string
	Attrs       map[string]string
	CreatedTime time.Time
}

var (
	UUID_NIL           gocql.UUID
	userPasswdMetadata = table.Metadata{
		Name:    "Person",
		Columns: []string{"namespace", "id", "pin", "email", "name", "first_name", "last_name", "other_names", "phones", "labels", "attrs", "created_time"},
		PartKey: []string{"namespace"},
		SortKey: []string{"id"},
	}
)
var personTable = table.New(userPasswdMetadata)

func CurrentNameSpace(ctx context.Context) string {
	return fmt.Sprint(ctx.Value(api.ContextPerson))
}

func (o *Person) Insert(ctx context.Context, session *gocqlx.Session) error {
	if o.Id == UUID_NIL {
		return fmt.Errorf("ID is empy")
	}
	o.Namespace = CurrentNameSpace(ctx)
	o.CreatedTime = time.Now()
	q := session.Query(personTable.Insert()).BindStruct(o)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("unable to insert loginPasswd(%v):%+v", o.Email, err)
	}
	return nil
}

// func GetByLabelKyes(ctx context.Context, session *gocqlx.Session, labels ...string) (*[]Person, error) {
func Get(ctx context.Context, session *gocqlx.Session) (*[]Person, error) {
	qb1 := qb.Select("person")
	filterValues := qb.M{}
	qb1 = AddFilterNS(ctx, qb1, &filterValues)
	return iter(qb1, &filterValues)
}

// func GetByLabelKyes(ctx context.Context, session *gocqlx.Session, labels ...string) (*[]Person, error) {
func GetByUUID(ctx context.Context, session *gocqlx.Session, uuid gocql.UUID) (*[]Person, error) {
	qb1 := qb.Select("person")
	qb1 = qb1.Where(qb.Eq("id"))
	filterValues := qb.M{}
	filterValues["id"] = uuid
	qb1 = AddFilterNS(ctx, qb1, &filterValues)
	return iter(qb1, &filterValues)
}

func AddFilterNS(ctx context.Context, parent *qb.SelectBuilder, filterValues *qb.M) *qb.SelectBuilder {
	out := parent.Where(qb.Eq("namespace"))
	(*filterValues)["namespace"] = CurrentNameSpace(ctx)
	return out
}

// func GetByLabelKyes(ctx context.Context, session *gocqlx.Session, labels ...string) (*[]Person, error) {
func GetByPin(ctx context.Context, session *gocqlx.Session, pin string) (*[]Person, error) {
	qb1 := qb.Select("person").Where(qb.Eq("pin"))
	log.Printf("qb: %v", qb1)
	filterValues := qb.M{}
	filterValues["pin"] = pin
	qb1 = AddFilterNS(ctx, qb1, &filterValues)
	return iter(qb1, &filterValues)
}

// func GetByLabelKyes(ctx context.Context, session *gocqlx.Session, labels ...string) (*[]Person, error) {
func GetByLabels(ctx context.Context, session *gocqlx.Session, labels ...string) (*[]Person, error) {
	labelCQL := "labels"

	qb1 := qb.Select("person").Where(qb.ContainsTuple(labelCQL, len(labels)))
	log.Printf("qb: %v", qb1)
	filterValues := qb.M{}
	qb1 = AddFilterNS(ctx, qb1, &filterValues)
	for k, v := range labels {
		filterName := fmt.Sprintf("%v_%d", labelCQL, k)
		filterValues[filterName] = v
		log.Printf("GetByLabels %v %v %v", k, v, filterName)
	}
	return iter(qb1, &filterValues)
}
