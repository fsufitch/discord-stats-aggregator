package analyzers

import discordstats "github.com/fsufitch/discord-stats-aggregator"

// UserMessageCountAnalyzer is an analyzer that tallies how many messages each user posted
type UserMessageCountAnalyzer struct {
	messageCounts map[string]int
}

type userMessageCountOutput struct {
	Counts map[string]int `json:"counts"`
}

// ID is used for creating the output structure
func (a UserMessageCountAnalyzer) ID() string {
	return "userMessageCounts"
}

// AddMessage implements the MessageRecipient interface
func (a *UserMessageCountAnalyzer) AddMessage(message *discordstats.CrawledMessage) error {
	if a.messageCounts == nil {
		a.messageCounts = map[string]int{}
	}
	username := message.Message.Author.Username
	if _, ok := a.messageCounts[username]; !ok {
		a.messageCounts[username] = 0
	}
	a.messageCounts[username]++
	return nil
}

// GetData implements the DataProvider interface
func (a UserMessageCountAnalyzer) GetData() (interface{}, error) {
	return userMessageCountOutput{a.messageCounts}, nil
}
