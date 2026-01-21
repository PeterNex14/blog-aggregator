# Blog Aggregator / RSS Feed Aggregator


**Gator** (short for aggreGATOR 🐊) is a CLI tool built in Go that functions as a personal RSS feed aggregator.

It allows users to:
-   **Add RSS feeds** from across the internet to be collected
-   **Store the collected posts** in a PostgreSQL database
-   **Follow and unfollow** RSS feeds that other users have added
-   **View summaries** of the aggregated posts in the terminal, with a link to the full post



## Prerequisites

To run this project, you need to have the following installed on your machine:

- **[Go](https://go.dev/doc/install)** (v1.23+): This is the programming language the application is written in. You need it to compile and run the code.
- **[PostgreSQL](https://www.postgresql.org/download/)**: This is the database used to store all the data. The application requires a running Postgres server with a database created.


Make sure your Postgres server is running and you have a connection string ready (usually in a config file or environment variable).

## Installation

To install the `gator` CLI tool, navigate to the root directory of the project and run:

```bash
go install .
```


This command compiles the project and places the executable binary in your `$GOPATH/bin` directory (typically `$HOME/go/bin` on Mac/Linux). Make sure this directory is in your system's `$PATH` so you can run the `gator` command from anywhere.

## Configuration

The application uses a JSON configuration file to store your database connection string and the current logged-in user.
The file is located at `~/.gatorconfig.json` (in your home directory).

Example configuration:

```json
{
  "db_url": "postgres://your_username:your_password@localhost:5432/your_database_name?sslmode=disable",
  "current_user_name": "your_username"
}
```

## Usage

Run the program using the `gator` command followed by a subcommand.

### Common Commands

*   **`register <name>`**: Register a new user.
*   **`login <name>`**: Login as an existing user.
*   **`users`**: List all registered users.
*   **`addfeed <name> <url>`**: Add a new RSS feed.
*   **`feeds`**: List all RSS feeds.
*   **`follow <url>`**: Follow an RSS feed.
*   **`browse <limit>`**: Browse posts from followed feeds (optional limit).
*   **`agg <time_between_requests>`**: Start the aggregator to fetch new posts (e.g., `1m`, `1h`).

Example:
```bash
gator register lane
gator login lane
gator addfeed "Boot.dev" "https://blog.boot.dev/index.xml"
gator agg 1m
```

## Project Background

This project serves as a practical application of backend development concepts, moving beyond simple scripts to a fully structured application.

### What I've Built

In this repository, I have implemented:
-   **A Persistent CLI**: A command-line interface that maintains state (user sessions, configuration) across executions.
-   **Database Integration**: A complete PostgreSQL schema with migrations to manage users, feeds, and feed follows.
-   **RSS Parsing**: Logic to fetch and parse XML data from remote RSS feeds.
-   **Worker Pools**: A background worker pattern to concurrently fetch updates from multiple feeds.
-   **Middleware**: Command-handler middleware design for authentication and access control.


### What I've Learned

Through building this project, I have mastered several key backend engineering concepts:

-   **Go & PostgreSQL Integration**: Learned how to integrate a Go application with a PostgreSQL database to create persistent, useful applications.
-   **Type-Safe SQL**: Practiced manually writing SQL schemas and queries, then leveraged tools like **sqlc** and **goose** to generate type-safe Go code. This ensures that database interactions are checked at compile-time, preventing runtime errors.
-   **Long-Running Services**: Implemented a long-running service (scraper) that continuously fetches new posts from RSS feeds in the background, dealing with concurrency and synchronization.

## Attribution

This project is part of the **[Boot.dev](https://boot.dev)** backend development curriculum. It is a guided project designed to teach "Backend Engineering" principles through hands-on practice.

## Extending the Project

You've done all the required steps, but if you'd like to make this project your own, here are some ideas:


- [ ] Add sorting and filtering options to the `browse` command
- [ ] Add pagination to the `browse` command
- [ ] Add concurrency to the `agg` command so that it can fetch more frequently
- [ ] Add a search command that allows for fuzzy searching of posts
- [ ] Add bookmarking or liking posts
- [ ] Add a TUI (Text User Interface) that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- [ ] Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- [ ] Write a service manager that keeps the `agg` command running in the background and restarts it if it crashes






