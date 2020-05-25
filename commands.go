package main

import (
	"fmt"
	log "github.com/inconshreveable/log15"
	"github.com/tidwall/redcon"
	"strconv"
	"strings"
	"errors"

)


func onRedconCommand(conn redcon.Conn, cmd redcon.Command) {
	redconCommandNext(conn, cmd)
}


func redconCommandNext(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToLower(string(cmd.Args[0]))

	switch command {
	case "ping":
		if len(cmd.Args) >= 2 {
			conn.WriteBulkString(string(cmd.Args[1]))
			return
		}
	case "add":

		if len(cmd.Args) < 4 || (len(cmd.Args)-2)%2 == 1 {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'add' command"))
			return
		}

		leaderboardName := string(cmd.Args[1])
		key := string(cmd.Args[2])
		score, err := strconv.ParseFloat(string(cmd.Args[3]), 64)

		if err != nil {
			writeRedconError(conn, err)
			return
		}

		fields := make(map[string]string)
		for i := 0; i < len(cmd.Args); i += 2 {
			if i >= 4 {
				fields[string(cmd.Args[i])] = string(cmd.Args[i+1])
			}
		}

		store.AddMember(leaderboardName,key,score,fields)
		conn.WriteString("OK")

	case "updatescore":
		if len(cmd.Args) < 3 {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'updatescore' command"))
		}

		leaderboardName := string(cmd.Args[1])
		key := string(cmd.Args[2])
		score, err := strconv.ParseFloat(string(cmd.Args[3]), 64)

		if err != nil {
			writeRedconError(conn, err)
			return
		}

		err = store.UpdateMemberScore(leaderboardName, key, score)

		if err != nil {
			writeRedconError(conn, err)
			return
		}

		conn.WriteString("OK")


	case "delmember":
		if len(cmd.Args) < 3  {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'delmember' command"))
		}
		leaderboardName := string(cmd.Args[1])
		key := string(cmd.Args[2])
		status := store.DeleteMember(leaderboardName, key)

		if status != nil {
			conn.WriteNull()
			return
		}
		conn.WriteString("OK")

	case "del":
		if len(cmd.Args) < 2  {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'del' command"))
		}
		leaderboardName := string(cmd.Args[1])
		status := store.DeleteLeaderboard(leaderboardName)

		if status != nil {
			conn.WriteNull()
			return
		}
		conn.WriteString("OK")


	case "score":
		if len(cmd.Args) < 3  {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'score' command"))
		}

		leaderboardName := string(cmd.Args[1])
		key := string(cmd.Args[2])

		score, err := store.MemberScore(leaderboardName,key)
		if err != nil {
			conn.WriteNull()
			return
		}

		scoreString := fmt.Sprintf("%.2f", score)
		conn.WriteString(scoreString)

	case "rank":
		if len(cmd.Args) < 3  {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'score' command"))
		}

		leaderboardName := string(cmd.Args[1])
		key := string(cmd.Args[2])

		rank, err := store.MemberRank(leaderboardName,key)
		if err != nil {
			conn.WriteNull()
			return
		}


		conn.WriteInt(rank)

	case "leaderboard":

		if len(cmd.Args) < 4 || (len(cmd.Args)-2)%2 == 1 {
			writeRedconError(conn, errors.New("ERR wrong number of arguments for 'leaderboard' command"))
			return
		}

		leaderboardName := string(cmd.Args[1])
		start, err := strconv.Atoi(string(cmd.Args[2]))

		if err != nil {
			writeRedconError(conn, errors.New("ERR start value must be integer"))
			return
		}

		if start < 0 {
			start = 0
		}

		stop, err := strconv.Atoi(string(cmd.Args[3]))

		if err != nil {
			writeRedconError(conn, errors.New("ERR stop value must be integer"))
			return
		}

		fields := make(map[string]string)
		for i := 0; i < len(cmd.Args); i += 2 {
			if i >= 4 {
				fields[string(cmd.Args[i])] = string(cmd.Args[i+1])
			}
		}

		result, err := store.GetSort(leaderboardName, fields)
		maxResult := len(result)

		if stop < 0 || stop > maxResult {
			stop = maxResult
		}

		if err != nil {
			conn.WriteNull()
			return
		}

		conn.WriteAny(result[start:stop])

	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	}
}

func writeRedconError(conn redcon.Conn, err error) {
	conn.WriteError(fmt.Sprintf("ERR on command: %s", err.Error()))
}

func onRedconConnect(conn redcon.Conn) bool {
	log.Info("Redcon new connection", "remote", conn.RemoteAddr())
	return true
}

func onRedconClose(conn redcon.Conn, err error) {
	log.Info("Redcon connection closed", "remote", conn.RemoteAddr())
}
