package chirps

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChirpWithLikes struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	User_id    int       `json:"user_id"`
	Created_at time.Time `json:"created_at"`
	Like_count int64     `json:"like_count"`
}

type Chirps struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	User_id    int       `json:"user_id"`
	Created_at time.Time `json:"created_at"`
}

type Repository interface {
	CreateChirp(ctx context.Context, content string, user_id int) error
	GetChirpById(ctx context.Context, id int) (*Chirps, error)
	GetChirpWithLikesById(ctx context.Context, id int) (*ChirpWithLikes, error)
	GetChirpsByUserId(ctx context.Context, user_id int) ([]Chirps, error)
	UpdateChirp(ctx context.Context, id int, content string) error
	DeleteChirp(ctx context.Context, id int) error
}

type dbChirpRepository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repository {
	return &dbChirpRepository{db: db}
}

func (dc *dbChirpRepository) CreateChirp(ctx context.Context, content string, user_id int) error {

	query := "INSERT INTO Chirps (content, user_id) VALUES ($1, $2)"

	_, err := dc.db.Exec(ctx, query, content, user_id)
	if err != nil {
		return fmt.Errorf("CreateChirp: could not insert new chirp: %w", err)
	}

	return nil
}

func (dc *dbChirpRepository) GetChirpWithLikesById(ctx context.Context, id int) (*ChirpWithLikes, error) {
	var chirp ChirpWithLikes

	query := "SELECT c.id, c.content, c.user_id, c.created_at, COUNT(cl.chirp_id) AS like_count FROM chirps c LEFT JOIN chirp_likes cl ON c.id=cl.chirp_id WHERE c.id=$1::int GROUP BY c.id"

	err := dc.db.QueryRow(ctx, query, id).Scan(
		&chirp.ID,
		&chirp.Content,
		&chirp.User_id,
		&chirp.Created_at,
		&chirp.Like_count,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("GetChirpWithLikesById: chirp with id %d not found", err)
		}

		return nil, fmt.Errorf("GetChirpWithLikesById: Error executing query %w", err)
	}

	return &chirp, nil

}

func (dc *dbChirpRepository) GetChirpById(ctx context.Context, id int) (*Chirps, error) {
	var chirp Chirps
	query := "SELECT * FROM Chirps WHERE id=$1"

	err := dc.db.QueryRow(ctx, query, id).Scan(
		&chirp.ID,
		&chirp.Content,
		&chirp.User_id,
		&chirp.Created_at,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("GetChirpById: chirp with id %d not found", err)
		}

		return nil, fmt.Errorf("GetChirpById: cant execute query %w", err)
	}

	return &chirp, nil

}

func (dc *dbChirpRepository) GetChirpsByUserId(ctx context.Context, user_id int) ([]Chirps, error) {

	var chirps []Chirps
	query := "SELECT * FROM Chirps WHERE user_id=$1"

	rows, err := dc.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetChirpsByUserId: cant execute query %w", err)
	}

	chirps, err = pgx.CollectRows(rows, pgx.RowToStructByName[Chirps])
	if err != nil {
		return nil, fmt.Errorf("GetChirpByUserId: cant getting rows %w", err)
	}

	return chirps, nil

}

func (dc *dbChirpRepository) UpdateChirp(ctx context.Context, id int, content string) error {
	var existingChirp Chirps

	selectQuery := "SELECT id, content FROM Chirps WHERE id=$1"

	err := dc.db.QueryRow(ctx, selectQuery, id, content).Scan(
		&existingChirp.ID,
		&existingChirp.Content,
		&existingChirp.User_id,
		&existingChirp.Created_at,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("UpdateChirp: chirp with id %d not found", id)
		}
		return fmt.Errorf("UpdateChirp: cant execute select query %w", err)
	}

	updateQuery := "UPDATE Chirps SET content=$1 WHERE id=$2"

	commandTag, err := dc.db.Exec(ctx, updateQuery, content, id)
	if err != nil {
		return fmt.Errorf("UpdateChirp: cant execute update query")
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("UpdateChirp: chirp with id %d not found", id)
	}

	return nil
}

func (dc *dbChirpRepository) DeleteChirp(ctx context.Context, id int) error {
	query := "DELETE FROM Chirps WHERE id=$1"

	commandTag, err := dc.db.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("DeleteChirp: could not execute delete query %w", err)
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("DeleteChirp: could not delete post: %w", err)
	}

	return nil
}

// func (dc *dbChirpRepository) EditChirp(ctx context.Context, chirp *Chirps) {}
