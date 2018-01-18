package filter

import discordstats "github.com/fsufitch/discord-stats-aggregator"

// ByBot creates a message filter for messages by/not by bots
func ByBot(bot bool) discordstats.MessageFilter {
	return func(m *discordstats.CrawledMessage) bool {
		return (bot && m.Message.Author.Bot) || (!bot && !m.Message.Author.Bot)
	}
}
