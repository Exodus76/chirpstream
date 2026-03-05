package userrelationship

import (
	"context"
	"fmt"
	"time"

	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/scylladb/gocqlx/v3/table"
)

// Table metadata for user relationships
var userFollowingTable = table.New(table.Metadata{
	Name:    "user_following",
	Columns: []string{"user_id", "following_id", "created_at"},
	PartKey: []string{"user_id"},
	SortKey: []string{"following_id"},
})

var userFollowersTable = table.New(table.Metadata{
	Name:    "user_followers",
	Columns: []string{"user_id", "follower_id", "created_at"},
	PartKey: []string{"user_id"},
	SortKey: []string{"follower_id"},
})

type Repository interface {
	FollowUser(ctx context.Context, userId int, followingId int) error
	UnfollowUser(ctx context.Context, userId int, followingId int) error
	GetFollowing(ctx context.Context, userId int) ([]int, error)
	GetFollowers(ctx context.Context, userId int) ([]int, error)
}

type dbUserRelationshipRepository struct {
	db *gocqlx.Session
}

func (d *dbUserRelationshipRepository) FollowUser(ctx context.Context, userId int, followingId int) error {

	stmt, names := userFollowingTable.Insert()
	err := d.db.Query(stmt, names).BindStruct(&struct {
		UserId      int       `db:"user_id"`
		FollowingId int       `db:"following_id"`
		CreatedAt   time.Time `db:"created_at"`
	}{
		UserId:      userId,
		FollowingId: followingId,
		CreatedAt:   time.Now(),
	}).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("FollowUser: Failed to insert into following %v", err)
	}

	stmt2, names2 := userFollowersTable.Insert()
	err = d.db.Query(stmt2, names2).BindStruct(&struct {
		UserId     int       `db:"user_id"`
		FollowerId int       `db:"follower_id"`
		CreatedAt  time.Time `db:"created_at"`
	}{
		UserId:     followingId,
		FollowerId: userId,
		CreatedAt:  time.Now(),
	}).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("FollowUser: Failed to insert into followers %v", err)
	}

	return nil

}

func (d *dbUserRelationshipRepository) GetFollowers(ctx context.Context, userId int) ([]int, error) {

	stmt, names := qb.Select("follower_id").Where(qb.Eq("user_id")).ToCql()
	iter := d.db.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id": userId,
	}).WithContext(ctx).Iter()

	var followers []int
	var followerId int
	for iter.Scan(&followerId) {
		followers = append(followers, followerId)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("GetFollowers: Failed to close iterator %v", err)
	}

	return followers, nil
}

func (d *dbUserRelationshipRepository) GetFollowing(ctx context.Context, userId int) ([]int, error) {

	stmt, names := qb.Select("following_id").Where(qb.Eq("user_id")).ToCql()
	iter := d.db.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id": userId,
	}).WithContext(ctx).Iter()

	var following []int
	var followingId int
	for iter.Scan(&followingId) {
		following = append(following, followingId)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("GetFollowing: Failed to close iterator %v", err)
	}

	return following, nil
}

func (d *dbUserRelationshipRepository) UnfollowUser(ctx context.Context, userId int, followingId int) error {
	panic("unimplemented")
}

func NewRepo(db *gocqlx.Session) Repository {
	return &dbUserRelationshipRepository{db: db}
}
