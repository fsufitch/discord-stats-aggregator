package filter

import (
	"time"

	"github.com/fsufitch/discord-stats-aggregator"
)

// AfterTime creates a message filter for messages after a time
func AfterTime(threshold time.Time) discordstats.MessageFilter {
	return func(m *discordstats.CrawledMessage) bool {
		msgTime, _ := m.Message.Timestamp.Parse()
		return msgTime.After(threshold)
	}
}

// BeforeTime creates a message filter for messages before a time
func BeforeTime(threshold time.Time) discordstats.MessageFilter {
	return func(m *discordstats.CrawledMessage) bool {
		msgTime, _ := m.Message.Timestamp.Parse()
		return msgTime.Before(threshold)
	}
}
