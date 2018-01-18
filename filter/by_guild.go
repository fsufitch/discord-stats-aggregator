package filter

import discordstats "github.com/fsufitch/discord-stats-aggregator"

// ByGuildID creates a message filter by several server IDs
func ByGuildID(guildIDs ...string) discordstats.MessageFilter {
	set := map[string]bool{}
	for _, id := range guildIDs {
		set[id] = true
	}

	return func(m *discordstats.CrawledMessage) bool {
		_, found := set[m.Channel.GuildID]
		return found
	}
}
