// Package cassandra holds the logic to interact with Cassandra DB
package cassandra

import (
	env "github.com/esequielvirtuoso/go_utils_lib/envs"
	"github.com/gocql/gocql"
)

const (
	envCassandraHost        = "OAUTH_CASSANDRA_HOST"
	envCassandraHostDefault = "127.0.0.1"
	envCassandraPort        = "OAUTH_CASSANDRA_PORT"
	envCassandraPortDefault = 9043
)

var (
	session *gocql.Session
)

func init() {
	// Connect to Cassandra cluster:
	cluster := gocql.NewCluster(getCassandraHost())
	cluster.Port = getCassandraPort()
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

func getCassandraPort() int {
	return env.GetInt(envCassandraPort, envCassandraPortDefault)
}

func getCassandraHost() string {
	return env.GetString(envCassandraHost, envCassandraHostDefault)
}
