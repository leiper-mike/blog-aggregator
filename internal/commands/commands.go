package commands

import (
	"fmt"
	"main/internal/config"
)

type State struct {
	Config *config.Config
}
type Command struct {
	Name string
	Args []string
}
type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: expected 1 args but recieved 0")
	}
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	println("User set as " + cmd.Args[0])
	return nil
}
func (c *Commands) Run(s *State, cmd Command) error {
	if command, ok := c.CommandMap[cmd.Name]; !ok {
		return fmt.Errorf("error: no command with name `%v` found in commands", cmd.Name)
	} else {
		err := command(s, cmd)
		return err
	}

}
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f
}
