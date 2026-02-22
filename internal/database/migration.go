package database

import (
	"fmt"
	"log"

	"github.com/scylladb/gocqlx/v3"
)

func Migrate(db *gocqlx.Session, keyspace string) error {
	log.Println("Starting database migration...")
	keyspaceQuery := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}", keyspace)

	statements := []string{
		keyspaceQuery,
		`
		CREATE TABLE IF NOT EXISTS chirps_by_user (
			user_id int,
			chirp_id timeuuid,
			content text,
			created_at timestamp,
			PRIMARY KEY ((user_id), chirp_id)
		) WITH CLUSTERING ORDER BY (chirp_id DESC);
	`,
		`
		CREATE TABLE IF NOT EXISTS chirp_stats(
			chirp_id timeuuid PRIMARY KEY,
			likes counter,
			retweets counter,
			replies counter
		);
	`,
	}

	for _, stmt := range statements {
		if err := db.Query(stmt, nil).Exec(); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("Database migration completed successfully!")
	return nil
}
