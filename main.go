package main

import (
	"fmt"
	"main/internal/commands"
	"main/internal/config"
	"os"
)

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("%v", err)
	}
	var state commands.State = commands.State{Config: &c}
	var coms commands.Commands = commands.Commands{CommandMap: make(map[string]func(*commands.State, commands.Command) error)}
	coms.Register("login", commands.HandlerLogin)

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
			fmt.Printf("%v", err)
			os.Exit(1)
		}
		os.Exit(0)

	}
}
