package database

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

func InitScyllaDb(host []string, username string, password string, keyspace string) (gocqlx.Session, error) {
	cluster := gocql.NewCluster(host...)

	// no need to verify
	cluster.IgnorePeerAddr = true
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}
	// cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("AWS_AP_SOUTH_1"))
	cluster.Keyspace = keyspace

	session, err := gocqlx.WrapSession(cluster.CreateSession())

	if err != nil {
		log.Printf("Failed connecting to scylladb")
		return session, err
	}

	return session, nil
}
