package main

import (
	"flag"
	log "github.com/inconshreveable/log15"
	"github.com/tidwall/redcon"
	"github.com/tinylib/msgp/msgp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagAddr = flag.String("addr", ":6488", "redcon port")
	flagDbPath = flag.String("dbPath", "leaderboard.db", "leaderboard db path")
	store = &Store{}
)

func main () {
	flag.Parse()

	leaderboardDB, err := os.OpenFile(*flagDbPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	store = NewStore()
	err = msgp.ReadFile(store, leaderboardDB)

	if err != nil {
		saveDb(store, leaderboardDB)
	}

	log.Info("LeaderboardDB Start", "port", *flagAddr)

	go runRedcon()
	go syncDb(store,leaderboardDB)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)



	<-c
	saveDb(store,leaderboardDB)

}


func runRedcon() {
	err := redcon.ListenAndServe(*flagAddr, onRedconCommand, onRedconConnect, onRedconClose)
	if err != nil {
		log.Crit("Redcon server startup failed", "error", err.Error())
		os.Exit(1)
	}
}

func syncDb(store *Store, db *os.File) {
	for {
		time.Sleep(2 * time.Minute)
		saveDb(store,db)
	}
}

func saveDb(store *Store, db *os.File) {
	msgp.WriteFile(store, db)
}