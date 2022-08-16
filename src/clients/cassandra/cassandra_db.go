package cassandra

import (
	"os"

	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func init() {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRACLUSTER"))
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}

}

func GetSession() *gocql.Session {
	return session
}
