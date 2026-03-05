package database

import (
	"fmt"
	"log"

	"github.com/scylladb/gocqlx/v3"
)

func MigrateChirp(db *gocqlx.Session, chirp_keyspace string) error {
	log.Println("Starting database migration...", chirp_keyspace)
	keyspaceQuery := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}", chirp_keyspace)

	if err := db.Query(keyspaceQuery, nil).Exec(); err != nil {
		return fmt.Errorf("Failed to create keyspace: %w", err)
	}

	statements := []string{
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

	log.Println("Chirp Database migration completed successfully!")
	return nil
}

func MigrateRelationship(db *gocqlx.Session, relationship_keyspace string) error {
	log.Println("Starting database migration...", relationship_keyspace)
	relationQuery := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}", relationship_keyspace)

	if err := db.Query(relationQuery, nil).Exec(); err != nil {
		return fmt.Errorf("Failed to create keyspace: %w", err)
	}

	statements := []string{
		relationQuery,
		`
		CREATE TABLE IF NOT EXISTS user_following (
			user_id int,
			following_id int,
			PRIMARY KEY ((user_id), following_id)
		);
	`,
		`
		CREATE TABLE IF NOT EXISTS user_followers (
			user_id int,
			follower_id int,
			PRIMARY KEY ((user_id), follower_id)
		);
	`,
	}

	for _, stmt := range statements {
		if err := db.Query(stmt, nil).Exec(); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("Relationship Database migration completed successfully!")
	return nil
}
