package discordstats

import "github.com/bwmarrin/discordgo"

// MessageFilter is a function for filtering Discord messages to process
type MessageFilter func(*CrawledMessage) bool

// MessageSource is an asynchronous source of Discord messages
type MessageSource interface {
	AddRecipients(recipients ...MessageRecipient)
	AddFilters(filters ...MessageFilter)

	StreamMessages() <-chan Progress
}

// Progress is a snapshot of the progress in trawling through a set of Discord messages
type Progress struct {
	MessagesRead     int
	MessagesRecorded int
	PercentComplete  float64
	Error            error
}

// CrawledMessage is a wrapper around a Discord message that is aware of its channel and previous/next messages
type CrawledMessage struct {
	Message *discordgo.Message
	Channel *discordgo.Channel
	Older   *CrawledMessage
	Newer   *CrawledMessage
}

// MessageRecipient is an object that can receive Discord messages and does something with them
type MessageRecipient interface {
	AddMessage(message *CrawledMessage) error
}

// DataProvider is an object that can provide an arbitrary serializable snapshot of some data
type DataProvider interface {
	GetData() (interface{}, error)
}

// MessageAnalyzer is a combined interface for receiving Discord messages and creating data
type MessageAnalyzer interface {
	ID() string
	MessageRecipient
	DataProvider
}
