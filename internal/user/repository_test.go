package user

import (
	"chirpstream/pkg/config"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

var testRepo Repository
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	conf, err := loadTestConfig()

	c := conf.Databases.Users_Test

	fmt.Println("Connecting to test database...")

	testPool, err = pgxpool.New(context.Background(), c.DBConnstring())
	if err != nil {
		log.Fatalf("Error while creating testing pool %v", err)
	}

	defer testPool.Close()

	testRepo = NewRepo(testPool)

	const migrationsDir = "../../internal/data/migrations"
	err = runMigrations(migrationsDir, testPool)
	if err != nil {
		log.Fatalf("Failed in run migrations: %v", err)
	}

	code := m.Run()

	os.Exit(code)

}

func TestRepo_handleCreateUser(t *testing.T) {
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE Users")

	//require is the same as assert but it stops execution when it fails
	require.NoError(t, err)

	newUser := &User{
		Name:       "Test",
		Email:      "test@test.com",
		User_name:  "mridul",
		Password:   "password",
		Active:     false,
		Created_at: time.Time{},
	}

	err = testRepo.CreateUser(context.Background(), newUser)

	//verify user is in
	var count int
	err = testPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM Users WHERE email=$1", "test@test.com").Scan(&count)

	require.NoError(t, err)
}

func runMigrations(dir string, testPool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(testPool)
	defer db.Close()

	//rollback to 0 adn then do the migration as we have to start from a CLEAN SLATE
	if err := goose.DownTo(db, dir, 0); err != nil {
		log.Printf("WARN: failed to run goose down (probably OK if DB is new): %v", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return err
	}
	return nil
}

// workaround for config not found issue
// TODO: find a better solution
func loadTestConfig() (config.Config, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	var err error
	for {
		if _, err := os.Stat(filepath.Join(basepath, "go.mod")); err == nil {
			break // Found the root
		}
		parent := filepath.Dir(basepath)
		if parent == basepath {
			return config.Config{}, err
		}
		basepath = parent
	}

	return config.Init(basepath)
}
