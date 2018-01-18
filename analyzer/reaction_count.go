package analyzer

import discordstats "github.com/fsufitch/discord-stats-aggregator"

// ReactionTally is an analyzer that tallies how many times each reaction was used
type ReactionTally struct {
	reactionCounts map[string]int
}

type reactionCountOutput struct {
	Counts map[string]int `json:"counts"`
}

// ID is used for creating the output structure
func (a ReactionTally) ID() string {
	return "reactionCounts"
}

// AddMessage implements the MessageRecipient interface
func (a *ReactionTally) AddMessage(message *discordstats.CrawledMessage) error {
	if a.reactionCounts == nil {
		a.reactionCounts = map[string]int{}
	}

	for _, reaction := range message.Message.Reactions {
		emoji := reaction.Emoji.Name
		if _, ok := a.reactionCounts[emoji]; !ok {
			a.reactionCounts[emoji] = 0
		}
		a.reactionCounts[emoji] += reaction.Count
	}
	return nil
}

// GetData implements the DataProvider interface
func (a ReactionTally) GetData() (interface{}, error) {
	return userMessageCountOutput{a.reactionCounts}, nil
}
