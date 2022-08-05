# caching-workshop
Workshop to accompany session about in-memory databases

## Running locally
First start a database using docker compose:
```bash
docker compose up
```

Then run the code with the 
```
go build; REDIS_PASSWORD=best_pass_ever FAST_RIVALS=true ./caching-workshop
```

## Interacting with the app
You can add some extra users with:
```bash
./data/insert_users.sh
```
And see the data with other commands like:
```bash
# Add extra users, optionally with starting points
curl -s  -H "Content-type: application/json" -d '{"username": "superstar", "points": 99}' localhost:8080/points
# Increment user score or add new user
curl -s  -H "Content-type: application/json" -d '{"username": "superstar"}' localhost:8080/points
# Get current values for eddie
curl -s -H "Content-type: application/json" localhost:8080/points/eddie
# See the entire leaderboard
curl -s -H "Content-type: application/json" localhost:8080/leaderboard
# Get the next best rival for eddie
curl -s -H "Content-type: application/json" localhost:8080/rival/eddie
```
