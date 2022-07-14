package db

import "time"

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
	// TODO implement faster response using cache
	return db.CalculateRivalSlow(username)
}
