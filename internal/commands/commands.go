package commands

import (
	"context"
	"fmt"
	"main/internal/config"
	"main/internal/database"
	"main/internal/rss"
	"strings"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db  *database.Queries
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
	name := cmd.Args[0]
	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return fmt.Errorf("user not found")
		} 
		return err
	}
	err = s.Config.SetUser(name)
	if err != nil {
		return err
	}
	println("User set as " + name)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: expected 1 args but recieved 0")
	}
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0] })
	if err != nil{
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User: %v was created, at %v, with ID: %v \n", user.Name,user.CreatedAt,user.ID)
	return nil
}

func HandlerReset(s *State, cmd Command) error{
	err := s.Db.DeleteAll(context.Background())
	if err != nil{
		return err
	}
	fmt.Println("Deleted all records in users table successfully")
	return nil
}
func HandlerListUsers(s *State, cmd Command) error{
	users, err := s.Db.GetAll(context.Background())
	if err != nil{
		return err
	}
	for _, user := range users{
		if user.Name == s.Config.CurrentUserName{
			fmt.Println(user.Name + " (current)")
		}else{
			fmt.Println(user.Name)
		}
	}
	return nil
	
}

func HandlerAgg(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil{
		return err
	}
	fmt.Println(feed.Format())
	for _, item := range feed.Channel.Item{
		fmt.Println(item.Format())
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("error: expected 2 args but recieved %v", len(cmd.Args))
	}
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil{
		return err
	}
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil{
		return err
	}
	fmt.Printf("ID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\nUrl: %v\nUserId: %v\n", feed.ID,feed.CreatedAt,feed.UpdatedAt,feed.Name,feed.Url,feed.UserID)
	return nil
}
func HandlerListFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds{
		user, err := s.Db.GetUserById(context.Background(),feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed Name: %v\nFeed URL: %v\nCreated By: %v\n", feed.Name,feed.Url, user.Name)
	}
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
