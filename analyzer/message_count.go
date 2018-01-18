package analyzer

import (
	"sort"

	discordstats "github.com/fsufitch/discord-stats-aggregator"
)

// UserMessageTally is an analyzer that tallies how many messages each user posted
type UserMessageTally struct {
	messageCounts map[string]int
}

// ID is used for creating the output structure
func (a UserMessageTally) ID() string {
	return "userMessageCounts"
}

// AddMessage implements the MessageRecipient interface
func (a *UserMessageTally) AddMessage(message *discordstats.CrawledMessage) error {
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

type userMessageTallyOutput []userMessageTallyOutputEntry
type userMessageTallyOutputEntry struct {
	Name  string `json:"username"`
	Count int    `json:"count"`
}

func (s userMessageTallyOutput) Len() int           { return len(s) }
func (s userMessageTallyOutput) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s userMessageTallyOutput) Less(i, j int) bool { return s[i].Count < s[j].Count }

// GetData implements the DataProvider interface
func (a UserMessageTally) GetData() (interface{}, error) {
	output := userMessageTallyOutput{}
	for username, count := range a.messageCounts {
		output = append(output, userMessageTallyOutputEntry{username, count})
	}
	sort.Sort(sort.Reverse(output))
	return output, nil
}
