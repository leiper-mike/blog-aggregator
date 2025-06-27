# About 
This is a RSS aggregator command line program that allows multiple users to subscribe to RSS feeds, aggregate, and browse posts from those feeds. 
# Setup
- [Go](https://go.dev/) and [PostgreSQL](https://www.postgresql.org/download/) are required to run the program.
- Navigate to the program folder, and run `go install`
- In your home directory, create a file named `.gatorconfig.json`
- Create a new or use an existing PostgreSQL database
- In the config file, fill in your connection string to your database
>`{
>  "db_url": "connection_string_goes_here",
>}`
- Using the schemas located in sql/schema, either use [goose](https://github.com/pressly/goose#install), or manually create the tables in your database.
- Build the program `go build`
# Usage
- Commands are structured like so: `main [command] [arguments]`
- `register [username]` registers a user
- `login [username]` logs in as a user
- `reset` resets the database
- `users` lists all users
- `addfeed [feedName] [feedUrl]` adds a feed to the database and subscribes the user to it
- `feeds` lists all feeds 
- `follow [feedUrl]` follows a feed (must already be in the database, use addFeed if it is not already in the database)
- `following` lists all feeds that the current user is subscribed to
- `unfollow [feedUrl]` unsubscribes a user to the feed
- `agg [duration]` Aggregates and displays all posts across all feeds that the user is subscribed to. Duration should be in the following format 1s, 2s, 1min, etc. and must be greater than 1 second to avoid straining the feed source.
- `browse [limit]` Browses all posts that were gathered using the `agg` command. Limit is an optional parameter.