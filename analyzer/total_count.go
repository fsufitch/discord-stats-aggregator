package analyzer

import discordstats "github.com/fsufitch/discord-stats-aggregator"

// TotalCount is an analyzer simply counts messages
type TotalCount struct {
	count int
}

// ID is used for creating the output structure
func (a TotalCount) ID() string {
	return "totalCount"
}

// AddMessage implements the MessageRecipient interface
func (a *TotalCount) AddMessage(message *discordstats.CrawledMessage) error {
	a.count++
	return nil
}

// GetData implements the DataProvider interface
func (a TotalCount) GetData() (interface{}, error) {
	return a.count, nil
}
