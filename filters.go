package discordstats

import "github.com/bwmarrin/discordgo"

// FilterByGuildID creates a message filter by several server IDs
func FilterByGuildID(guildIDs ...string) MessageFilter {
	set := map[string]bool{}
	for _, id := range guildIDs {
		set[id] = true
	}

	return func(m *discordgo.Message, c *discordgo.Channel) bool {
		_, found := set[c.GuildID]
		return found
	}
}

// FilterByPublicChannel creates a message filter for channels that are/are not public
func FilterByPublicChannel(public bool) MessageFilter {
	return func(m *discordgo.Message, c *discordgo.Channel) bool {
		return (public && c.Type == discordgo.ChannelTypeGuildText) ||
			(!public && c.Type == discordgo.ChannelTypeDM)
	}
}

// FilterByBot creates a message filter for messages by/not by bots
func FilterByBot(bot bool) MessageFilter {
	return func(m *discordgo.Message, c *discordgo.Channel) bool {
		return (bot && m.Author.Bot) || (!bot && !m.Author.Bot)
	}
}
