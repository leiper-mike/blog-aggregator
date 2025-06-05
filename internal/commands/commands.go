package commands

import (
	"context"
	"database/sql"
	"fmt"
	"main/internal/config"
	"main/internal/database"
	"main/internal/rss"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
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
		if strings.Contains(err.Error(), "no rows in result set") {
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
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{ID: int32(uuid.New().ID()), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]})
	if err != nil {
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User: %v was created, at %v, with ID: %v \n", user.Name, user.CreatedAt, user.ID)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAll(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Deleted all records in users table successfully")
	return nil
}
func HandlerListUsers(s *State, cmd Command) error {
	users, err := s.Db.GetAll(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil

}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: expected 1 args but recieved 0")
	}
	dur := cmd.Args[0]
	interval, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}
	oneS, _ := time.ParseDuration("1s")
	if interval < oneS {
		return fmt.Errorf("error: interval cannot be less than one second")
	}
	fmt.Println("Collecting feeds every " + interval.String())
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		dbFeed, err := scrapeFeeds(s)
		if err != nil {
			return err
		}
		url := dbFeed.Url
		feed, err := rss.FetchFeed(context.Background(), url)
		if err != nil {
			return err
		}
		fmt.Println(feed.Format())
		for _, item := range feed.Channel.Item {
			fmt.Println(item.Format())
			pub, err := parseTime(item.PubDate)
			isValid := true
			if err != nil{
				isValid = false
			}
			_, err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
				ID:          int32(uuid.New().ID()),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: true},
				PublishedAt: sql.NullTime{Time: pub, Valid: isValid},
				FeedID:      dbFeed.ID,
			})
			
			if err != nil {
				if !strings.Contains(err.Error(), "duplicate key"){ 
					return err
				} 	
				
			}
		}
	}
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("error: expected 2 args but recieved %v", len(cmd.Args))
	}
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("ID: %v\nCreatedAt: %v\nUpdatedAt: %v\nName: %v\nUrl: %v\nUserId: %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
func HandlerBrowse(s *State, cmd Command, user database.User) error {
	limit := 2
	if len(cmd.Args) != 0 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err != nil{
			l= 2
		}
		limit = l
	}
	posts, err := s.Db.GetPostsForUser(context.Background(),database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	for _, post := range posts{
		feed, err := s.Db.GetFeedByID(context.Background(), post.FeedID)
		if err != nil{
			return err
		}
		fmt.Printf("Title: %v\nPublished: %v\nDescription: %v\nURL: %v\nFeed: %v\n", post.Title, post.PublishedAt.Time, post.Description.String, post.Url, feed.Name)
		
	}
	return nil
}
func HandlerListFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed Name: %v\nFeed URL: %v\nCreated By: %v\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}
func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: expected 1 args but recieved 0")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("User %v is now following feed '%v'\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}
func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsByUserId(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%v is following: \n", s.Config.CurrentUserName)
	for i, feedFollow := range feedFollows {
		fmt.Printf("%v: %v\n", i, feedFollow.FeedName)
	}
	return nil
}
func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("error: expected 1 args but recieved 0")
	}
	url := cmd.Args[0]
	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v unfollowed feed %v with url: %v", user.Name, feed.Name, feed.Url)
	return err
}
func scrapeFeeds(s *State) (database.Feed, error) {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return feed, err
	}
	err = s.Db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	return feed, err
}
func parseTime(date string) (time.Time, error) {
	t, err := time.Parse(time.UnixDate, date)
	if err != nil {
		t, err = time.Parse(time.DateTime, date)
	}
	if err != nil {
		t, err = time.Parse(time.DateOnly, date)
	}
	if err != nil {
		t, err = time.Parse(time.ANSIC, date)
	}
	if err != nil {
		t, err = time.Parse(time.RubyDate, date)
	}
	if err != nil {
		t, err = time.Parse(time.Stamp, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC850, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC1123, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC1123Z, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC822, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC822Z, date)
	}
	if err != nil {
		t, err = time.Parse(time.RFC3339, date)
	}
	if err != nil {
		t, err = time.Parse(time.Kitchen, date)
	}
	return t, err
}
func (c *Commands) Run(s *State, cmd Command) error {
	if command, ok := c.CommandMap[cmd.Name]; !ok {
		return fmt.Errorf("error: no command with name `%v` found in commands", cmd.Name)
	} else {
		err := command(s, cmd)
		return err
	}

}
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f
}
