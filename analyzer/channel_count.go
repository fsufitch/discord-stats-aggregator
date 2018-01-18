package analyzer

import (
	"sort"

	discordstats "github.com/fsufitch/discord-stats-aggregator"
)

// ChannelTally is an analyzer that tallies how many messages each user posted
type ChannelTally struct {
	messageCounts map[string]int
}

// ID is used for creating the output structure
func (a ChannelTally) ID() string {
	return "channelMessageCounts"
}

// AddMessage implements the MessageRecipient interface
func (a *ChannelTally) AddMessage(message *discordstats.CrawledMessage) error {
	if a.messageCounts == nil {
		a.messageCounts = map[string]int{}
	}
	channel := message.Channel.Name
	if _, ok := a.messageCounts[channel]; !ok {
		a.messageCounts[channel] = 0
	}
	a.messageCounts[channel]++
	return nil
}

type channelTallyOutput []channelTallyOutputEntry
type channelTallyOutputEntry struct {
	ChannelName string `json:"channelName"`
	Count       int    `json:"count"`
}

func (s channelTallyOutput) Len() int           { return len(s) }
func (s channelTallyOutput) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s channelTallyOutput) Less(i, j int) bool { return s[i].Count < s[j].Count }

// GetData implements the DataProvider interface
func (a ChannelTally) GetData() (interface{}, error) {
	output := channelTallyOutput{}
	for channelName, count := range a.messageCounts {
		output = append(output, channelTallyOutputEntry{channelName, count})
	}
	sort.Sort(sort.Reverse(output))
	return output, nil
}
