package db

import (
	"github.com/go-redis/redis/v8"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Points   int    `json:"points"`
	Rank     int    `json:"rank"`
}

func (db *Database) SaveUser(user *User) (*User, error) {
	pipe := db.Client.TxPipeline()

	if user.Points != 0 {
		member := &redis.Z{
			Score:  float64(user.Points),
			Member: user.Username,
		}
		pipe.ZAdd(Ctx, leaderboardKey, member)
	} else {
		pipe.ZIncrBy(Ctx, leaderboardKey, 1, user.Username)
	}

	rank := pipe.ZRevRank(Ctx, leaderboardKey, user.Username)
	points := pipe.ZScore(Ctx, leaderboardKey, user.Username)
	_, err := pipe.Exec(Ctx)

	if err != nil {
		return nil, ErrNil
	}
	user.Rank = int(rank.Val())
	user.Points = int(points.Val())
	return user, nil
}

func (db *Database) GetUser(username string) (*User, error) {
	pipe := db.Client.TxPipeline()
	score := pipe.ZScore(Ctx, leaderboardKey, username)
	rank := pipe.ZRevRank(Ctx, leaderboardKey, username)
	_, err := pipe.Exec(Ctx)
	if err != nil {
		return nil, err
	}
	if score == nil {
		return nil, ErrNil
	}
	return &User{
		Username: username,
		Points:   int(score.Val()),
		Rank:     int(rank.Val()),
	}, nil
}
