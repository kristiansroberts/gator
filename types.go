package main

import (
	"fmt"

	"github.com/kristiansroberts/gator/internal/config"
	"github.com/kristiansroberts/gator/internal/database"
)

type state struct {
	db     *database.Queries
	cfgPtr *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *commands) run(state *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(state, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.registeredCommands[name] = handler
}
