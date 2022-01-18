package cassandra

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ggrrrr/bui_lib/config"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/spf13/viper"
)

const (
	CQL_CLUSTER         = "cql.cluster"
	CQL_KEYSPACE        = "cql.keyspace.name"
	CQL_CREATE_KEYSPACE = "cql.keyspace.replication"
	CQL_FILES_BASEDIR   = "cql.files.basedir"
)

var (
	envParamsDefaults = []config.ParamValue{
		{
			Name:     CQL_CLUSTER,
			Info:     "Cassandra/Scylla cluster connection <HOST1>:<PORT1>,<HOST2>:<PORT2>,",
			DefValue: "127.0.0.1",
		},
		{
			Name:     CQL_KEYSPACE,
			Info:     "Cassandra/Scylla keyspace name",
			DefValue: "test",
		},
		{
			Name:     CQL_CREATE_KEYSPACE,
			Info:     "create keyspace ... with replication = {'class': 'SimpleStrategy','replication_factor' : 1}, replication = {'class': 'SimpleStrategy','replication_factor' : 3}",
			DefValue: "replication = {'class': 'SimpleStrategy','replication_factor' : 1}",
		},
		{
			Name:     CQL_FILES_BASEDIR,
			Info:     "Base directory with all CQL for createing schemas.",
			DefValue: "/app/cql",
		},
	}

	// dbCluster  []string
	cqlBaseDir string
	dbKeyspace string
	Session    *gocqlx.Session
	cluster    *gocql.ClusterConfig
)

type SomeType struct{}

func (o *SomeType) ObserveConnect(c gocql.ObservedConnect) {
	// log.Printf("ObserveConnect: %+v", c)
}

func Configure() error {
	config.Configure(envParamsDefaults)
	hosts := viper.GetString(CQL_CLUSTER)
	dbKeyspace = viper.GetString(CQL_KEYSPACE)
	qCreateKeyspaceWith := viper.GetString(CQL_CREATE_KEYSPACE)
	cqlBaseDir = viper.GetString(CQL_FILES_BASEDIR)
	if hosts == "" {
		fmt.Println(config.Help())
		return fmt.Errorf("param: %s not set", CQL_CLUSTER)
	}
	hostss := strings.Split(hosts, ",")
	if len(hostss) == 0 {
		fmt.Println(config.Help())
		return config.ErrorParamInvalid(CQL_CLUSTER, hosts)
	}
	// dbCluster = hostss
	if dbKeyspace == "" {
		fmt.Println(config.Help())
		return fmt.Errorf("param: %s not set", CQL_KEYSPACE)
	}
	log.Printf("Connecting:%v, keyspace:%v", hostss, dbKeyspace)

	cluster = gocql.NewCluster(hostss...)
	cluster.Timeout = 10 * time.Second

	obCon := &SomeType{}
	// obCon := &SomeType{}
	cluster.ConnectObserver = obCon

	return checkKeyspace(dbKeyspace, qCreateKeyspaceWith)
}

func Shutdown() {
	if Session != nil {
		Session.Close()
		log.Printf("session.closed")
	}
}

func Connect() (*gocqlx.Session, error) {
	// cluster.ConnectObserver.ObserveConnect()
	// defer cluster.
	cluster.Keyspace = dbKeyspace
	s, err := gocqlx.WrapSession(cluster.CreateSession())

	// s, err := cluster.CreateSession()
	if err != nil {
		// panic(fmt.Sprintf("db error: %+v", err))
		return nil, err
	}
	log.Printf("Connected: keyspace: %s @  [%v]", dbKeyspace, cluster.Hosts)

	Session = &s
	// testQ()
	return &s, nil
	// defer Session.Close()
	// err = Session.Query("CREATE KEYSPACE pesho WITH replication = {'class': 'SimpleStrategy','replication_factor':1}").
	// 	Exec()
	// log.Printf("%+v", err)

}
