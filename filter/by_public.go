package filter

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fsufitch/discord-stats-aggregator"
)

// ByPublicChannel creates a message filter for channels that are/are not public
func ByPublicChannel(public bool) discordstats.MessageFilter {
	return func(m *discordstats.CrawledMessage) bool {
		return (public && m.Channel.Type == discordgo.ChannelTypeGuildText) ||
			(!public && m.Channel.Type == discordgo.ChannelTypeDM)
	}
}
