package analyzer

import (
	"sort"

	discordstats "github.com/fsufitch/discord-stats-aggregator"
)

// ReactionTally is an analyzer that tallies how many times each reaction was used
type ReactionTally struct {
	reactionCounts map[string]int
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

type reactionTallyOutput []reactionTallyOutputEntry
type reactionTallyOutputEntry struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

func (s reactionTallyOutput) Len() int           { return len(s) }
func (s reactionTallyOutput) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s reactionTallyOutput) Less(i, j int) bool { return s[i].Count < s[j].Count }

// GetData implements the DataProvider interface
func (a ReactionTally) GetData() (interface{}, error) {
	output := reactionTallyOutput{}
	for emoji, count := range a.reactionCounts {
		output = append(output, reactionTallyOutputEntry{emoji, count})
	}
	sort.Sort(sort.Reverse(output))
	return output, nil
}
