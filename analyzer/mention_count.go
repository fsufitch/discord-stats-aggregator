package analyzer

import (
	"regexp"
	"sort"

	discordstats "github.com/fsufitch/discord-stats-aggregator"
)

// HereMentionTally is an analyzer that talllies who used the most @here
type HereMentionTally struct {
	mentionCounts map[string]int
}

// ID is used for creating the output structure
func (a HereMentionTally) ID() string {
	return "@hereMentionTally"
}

var hereMentionRegexp = regexp.MustCompile("(^|\\W)@here\\W")

// AddMessage implements the MessageRecipient interface
func (a *HereMentionTally) AddMessage(message *discordstats.CrawledMessage) error {
	if a.mentionCounts == nil {
		a.mentionCounts = map[string]int{}
	}

	if !hereMentionRegexp.MatchString(message.Message.Content) {
		return nil
	}

	if _, ok := a.mentionCounts[message.Message.Author.Username]; !ok {
		a.mentionCounts[message.Message.Author.Username] = 0
	}
	a.mentionCounts[message.Message.Author.Username]++
	return nil
}

type mentionTallyOutput []mentionTallyOutputEntry
type mentionTallyOutputEntry struct {
	Name  string
	Count int
}

func (s mentionTallyOutput) Len() int           { return len(s) }
func (s mentionTallyOutput) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s mentionTallyOutput) Less(i, j int) bool { return s[i].Count < s[j].Count }

// GetData implements the DataProvider interface
func (a HereMentionTally) GetData() (interface{}, error) {
	output := mentionTallyOutput{}
	for username, mentionCount := range a.mentionCounts {
		output = append(output, mentionTallyOutputEntry{username, mentionCount})
	}
	sort.Sort(sort.Reverse(output))

	return output, nil
}
