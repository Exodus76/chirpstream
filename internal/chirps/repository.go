package chirps

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
)

// Define table metadata once
var chirpTable = table.New(table.Metadata{
	Name:    "chirps_by_user",
	Columns: []string{"user_id", "chirp_id", "content", "created_at"},
	PartKey: []string{"user_id"},
	SortKey: []string{"chirp_id"},
})

var statsTable = table.New(table.Metadata{
	Name:    "chirp_stats",
	Columns: []string{"chirp_id", "likes", "retweets", "replies"},
	PartKey: []string{"chirp_id"},
})

type Repository interface {
	CreateChirp(ctx context.Context, content string, userId int) error
	GetChirpById(ctx context.Context, userId int, chirpId gocql.UUID) (*Chirp, error)
	GetChirpsByUserId(ctx context.Context, userId int, pageState []byte, limit int) ([]Chirp, []byte, error)
	UpdateChirp(ctx context.Context, userId int, chirpId gocql.UUID, content string) error
	DeleteChirp(ctx context.Context, userId int, chirpId gocql.UUID) error
}

type dbChirpRepository struct {
	db *gocqlx.Session
}

func NewRepo(db *gocqlx.Session) Repository {
	return &dbChirpRepository{db: db}
}

// TODO: need to add updatedAt
type Chirp struct {
	ChirpId   gocql.UUID `db:"chirp_id"`
	UserId    int        `db:"user_id"`
	Name      string     `db:"name"`
	Username  string     `db:"username"`
	Content   string     `db:"content"`
	Likes     int        `db:"likes"`
	Retweets  int        `db:"retweets"`
	Replies   int        `db:"replies"`
	CreatedAt time.Time  `db:"created_at"`
}

func (dc *dbChirpRepository) CreateChirp(ctx context.Context, content string, userId int) error {
	chirp := Chirp{
		UserId:    userId,
		ChirpId:   gocql.TimeUUID(),
		Content:   content,
		CreatedAt: time.Now(),
	}

	stmt, names := chirpTable.Insert()

	// BindStruct maps the Go struct fields directly to the query
	return dc.db.Query(stmt, names).BindStruct(chirp).Exec()
}

func (dc *dbChirpRepository) GetChirpById(ctx context.Context, userId int, chirpId gocql.UUID) (*Chirp, error) {
	var chirp Chirp

	stmt, names := chirpTable.Get("user_id", "chirp_id")
	err := dc.db.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id":  userId,
		"chirp_id": chirpId,
	}).Get(&chirp)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("GetChirpById: chirp with id %s not found", chirpId)
		}
		return nil, fmt.Errorf("GetChirpById: cant execute query %w", err)
	}

	stmtStats, namesStats := statsTable.Get("chirp_id")
	err = dc.db.Query(stmtStats, namesStats).BindMap(map[string]interface{}{
		"chirp_id": chirpId,
	}).Get(&chirp)

	if err != nil && err != gocql.ErrNotFound {
		return nil, fmt.Errorf("GetChirpById: cant execute stats query %w", err)
	}

	return &chirp, nil
}

func (dc *dbChirpRepository) GetChirpsByUserId(ctx context.Context, userId int, pageState []byte, limit int) ([]Chirp, []byte, error) {
	var chirps []Chirp

	// stmt, names := chirpTable.Get("user_id")
	//.Get is for single row, we need to use Select for multiple rows
	stmt, names := qb.Select(chirpTable.Name()).Where(qb.Eq("user_id")).ToCql()
	iter := dc.db.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id": userId,
	}).PageSize(limit).PageState(pageState).Iter()

	err := iter.Select(&chirps)
	nextPageState := iter.PageState()

	if closeErr := iter.Close(); closeErr != nil {
		return nil, nil, fmt.Errorf("cant execute query: %w", closeErr)
	}

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil, fmt.Errorf("GetChirpsByUserId: no chirps found for user id %d", userId)
		}
		return nil, nil, fmt.Errorf("GetChirpsByUserId: cant execute query %w", err)
	}

	return chirps, nextPageState, nil
}

func (dc *dbChirpRepository) UpdateChirp(ctx context.Context, userId int, chirpId gocql.UUID, content string) error {
	stmt, names := chirpTable.Update("content")

	err := dc.db.Query(stmt, names).BindMap(map[string]interface{}{
		"content":  content,
		"user_id":  userId,
		"chirp_id": chirpId,
	}).Exec()

	if err != nil {
		return fmt.Errorf("UpdateChirp: cant execute update query %w", err)
	}
	return nil
}

func (dc *dbChirpRepository) DeleteChirp(ctx context.Context, userId int, chirpId gocql.UUID) error {
	stmt, names := chirpTable.Delete()

	err := dc.db.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id":  userId,
		"chirp_id": chirpId,
	}).Exec()

	if err != nil {
		return fmt.Errorf("DeleteChirp: could not execute delete query %w", err)
	}
	return nil
}
