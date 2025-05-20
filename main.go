package main

import (
	"database/sql"
	"fmt"
	"main/internal/commands"
	"main/internal/config"
	"main/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}
	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil{
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	var state commands.State = commands.State{Db: dbQueries, Config: &c}
	var coms commands.Commands = commands.Commands{CommandMap: make(map[string]func(*commands.State, commands.Command) error)}
	coms.Register("login", commands.HandlerLogin)
	coms.Register("register", commands.HandlerRegister)
	coms.Register("reset", commands.HandlerReset)
	coms.Register("users", commands.HandlerListUsers)
	coms.Register("agg", commands.HandlerAgg)
	coms.Register("addfeed", commands.HandlerAddFeed)
	coms.Register("feeds", commands.HandlerListFeeds)
	coms.Register("follow", commands.HandlerFollow)
	coms.Register("following", commands.HandlerFollowing)
	a := os.Args
	if len(a) < 2 {
		fmt.Printf("error: no command provided \n")
		os.Exit(1)
	} else {
		cmdName := a[1]
		var args []string
		if len(a) > 2 {
			args = a[2:]
		}
		err := coms.Run(&state, commands.Command{Name: cmdName, Args: args})
		if err != nil {
			fmt.Printf("%v \n", err)
			os.Exit(1)
		}
		os.Exit(0)

	}
}
