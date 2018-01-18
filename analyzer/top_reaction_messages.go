package analyzer

import (
	"sort"

	discordstats "github.com/fsufitch/discord-stats-aggregator"
)

// TopReactionMessages is an analyzer that shows the best messages for each reaction
type TopReactionMessages struct {
	reactionMessages map[string][]reactionEntry
}

type reactionEntry struct {
	ReactionCount int
	Message       *discordstats.CrawledMessage
}

type topReactionMessagesOutput map[string][]topReactionOutputEntry
type topReactionOutputEntry struct {
	ReactionCount   int    `json:"reactionCount"`
	User            string `json:"user"`
	Channel         string `json:"channel"`
	MessageContents string `json:"message"`
	TimestampUTC    string `json:"timestamp"`
}

// ID is used for creating the output structure
func (a TopReactionMessages) ID() string {
	return "topReactionMessages"
}

// AddMessage implements the MessageRecipient interface
func (a *TopReactionMessages) AddMessage(message *discordstats.CrawledMessage) error {
	if a.reactionMessages == nil {
		a.reactionMessages = map[string][]reactionEntry{}
	}

	for _, reaction := range message.Message.Reactions {
		emoji := reaction.Emoji.Name
		a.reactionMessages[emoji] = append(a.reactionMessages[emoji], reactionEntry{reaction.Count, message})
		a.reactionMessages[emoji] = topReactionMessages(a.reactionMessages[emoji], 3)
	}
	return nil
}

type reactionEntrySorter []reactionEntry

func (s reactionEntrySorter) Len() int           { return len(s) }
func (s reactionEntrySorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s reactionEntrySorter) Less(i, j int) bool { return s[i].ReactionCount < s[j].ReactionCount }

func topReactionMessages(entries []reactionEntry, topCount int) []reactionEntry {
	copyMessages := append([]reactionEntry{}, entries...)
	sort.Sort(sort.Reverse(reactionEntrySorter(copyMessages)))
	if len(copyMessages) > topCount {
		copyMessages = copyMessages[0:topCount]
	}
	return copyMessages
}

// GetData implements the DataProvider interface
func (a TopReactionMessages) GetData() (interface{}, error) {
	output := topReactionMessagesOutput{}
	for emoji, reactionEntries := range a.reactionMessages {
		output[emoji] = []topReactionOutputEntry{}
		for _, entry := range reactionEntries {
			timestamp, _ := entry.Message.Message.Timestamp.Parse()
			outputEntry := topReactionOutputEntry{
				ReactionCount:   entry.ReactionCount,
				User:            entry.Message.Message.Author.Username,
				Channel:         entry.Message.Channel.Name,
				MessageContents: entry.Message.Message.ContentWithMentionsReplaced(),
				TimestampUTC:    timestamp.UTC().String(),
			}
			output[emoji] = append(output[emoji], outputEntry)
		}
	}
	return output, nil
}
