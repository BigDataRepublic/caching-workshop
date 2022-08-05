package db

import (
	"log"
	"time"
)

var rivalKey = "rival"

func (db *Database) CalculateRivalSlow(username string) (*User, error) {
	rank := db.Client.ZRevRank(Ctx, leaderboardKey, username).Val()
	time.Sleep(2 * time.Second)
	if rank == 0 {
		rivalData := db.Client.ZRevRangeWithScores(Ctx, leaderboardKey, -1, -1).Val()[0]
		return &User{
			Username: rivalData.Member.(string),
			Points:   int(rivalData.Score),
			Rank:     1,
		}, nil
	} else {
		rivalData := db.Client.ZRevRangeWithScores(Ctx, leaderboardKey, rank-1, rank).Val()[0]
		return &User{
			Username: rivalData.Member.(string),
			Points:   int(rivalData.Score),
			Rank:     int(rank - 1),
		}, nil
	}

}

func (db *Database) GetRival(username string) (*User, error) {
	rival := db.Client.HGet(Ctx, rivalKey, username).Val()
	if len(rival) == 0 {
		log.Printf("Rival data not available, will calculate ad hoc")
		return db.CalculateRivalSlow(username)
	} else {
		return db.GetUser(rival)
	}
}

func (db *Database) UpdateRivals() error {
	rivalData := db.Client.ZRevRange(Ctx, leaderboardKey, 0, -1).Val()
	mappedData := make(map[string]string)
	for i := 0; i < len(rivalData); i += 1 {
		// First player will just need to fight the second always
		rivalDataPlus := append([]string{rivalData[1]}, rivalData...)
		mappedData[rivalData[i]] = rivalDataPlus[i]
	}
	// This is really difficult
	time.Sleep(2 * time.Second)

	return db.Client.HSet(Ctx, rivalKey, mappedData).Err()
}
