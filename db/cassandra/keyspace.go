package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

const Q_CREATE_KEYSPACE = "CREATE KEYSPACE IF NOT EXISTS %s WITH \n\t%s"

const Q_CREATE_T_V = "CREATE TABLE apiVersion ( c text PRIMARY KEY (c) )"

func checkKeyspace(keyspace string, qCreateKeyspaceWith string) error {
	c := gocql.NewCluster(cluster.Hosts...)

	c.Keyspace = "system"
	session, err := c.CreateSession()
	if err != nil {
		log.Fatal("createSession:", err)
	}
	defer session.Close()
	q := fmt.Sprintf(Q_CREATE_KEYSPACE, keyspace, qCreateKeyspaceWith)
	log.Printf("db: %s", q)
	return session.Query(q).Exec()

}
