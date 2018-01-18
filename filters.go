package discordstats

import "github.com/bwmarrin/discordgo"

// FilterByGuildID creates a message filter by several server IDs
func FilterByGuildID(guildIDs ...string) MessageFilter {
	set := map[string]bool{}
	for _, id := range guildIDs {
		set[id] = true
	}

	return func(m *CrawledMessage) bool {
		_, found := set[m.Channel.GuildID]
		return found
	}
}

// FilterByPublicChannel creates a message filter for channels that are/are not public
func FilterByPublicChannel(public bool) MessageFilter {
	return func(m *CrawledMessage) bool {
		return (public && m.Channel.Type == discordgo.ChannelTypeGuildText) ||
			(!public && m.Channel.Type == discordgo.ChannelTypeDM)
	}
}

// FilterByBot creates a message filter for messages by/not by bots
func FilterByBot(bot bool) MessageFilter {
	return func(m *CrawledMessage) bool {
		return (bot && m.Message.Author.Bot) || (!bot && !m.Message.Author.Bot)
	}
}
