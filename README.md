# LeaderboardDB
> LeaderboardDB is storage for sorted lists with custom field filters. Leaderboard uses the Redis RESP protocol.

![LeaderboardDB Logo](https://github.com/fatihkahveci/leaderboard-db/300x300.png)

## Installation
```
go get -u github.com/fatihkahveci/leaderboard-db
```
## Flags

| Flag        | Desc           | 
| ------------- |-------------| 
| addr | LeaderboardDB port --> Default: :6488 | 
| dbPath | LeaderboardDB dbpath --> Default: leaderboard.db | 

## Commands 

| Command | Params | Example |
| ------ | ------ |----------- |
| add   | leaderboardKey, memberKey, score, [fields] (optional) | add diablo_3 user_1 123 character demon_hunter |
| leaderboard | leaderboardKey, start, stop, [fields] (optional) | leaderboard diablo_3 0 -1 character demon_hunter |
| del    | leaderboardKey | del diablo_3 |
| delmember    | leaderboardKey, key | del diablo_3 user_1 |
| updatescore    | leaderboardKey, key, score | updatescore diablo_3 user_1 3 |

## Examples
Let's say you need to store Diablo 3 player scores and need to filter character and country.

First we need to add users to diablo_3 leaderboard.
```
> add diablo_3 user_1 100 character wizard country tr
> add diablo_3 user_2 105 character monk country tr
> add diablo_3 user_3 95 character monk country tr
> add diablo_3 user_4 95 character monk country us
> add diablo_3 user_5 96 character demon_hunter country us
```

If you need to get all users in tihs leaderboard playing monk.
```
leaderboard diablo_3 0 -1 character monk
```
Response will be:

```
1) "user_2"
2) "user_3"
3) "user_4"
```

Or maybe you need to get all monk player in some country.

```
leaderboard diablo_3 0 -1 character monk country tr
```

Then response will be:

```
1) "user_2"
2) "user_3"
```

Or maybe you don't need any filter

```
leaderboard diablo_3 0 -1
```

Then response will be:

```
1) "user_2"
2) "user_1"
3) "user_5"
4) "user_3"
5) "user_4"
```

## Thanks

Thanks for logo @eraydemirok
