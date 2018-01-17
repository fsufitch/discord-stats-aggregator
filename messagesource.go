package discordstats

import "github.com/bwmarrin/discordgo"

// BasicMessageSource is a MessageSource based on a Discord session and a time bracket
type BasicMessageSource struct {
	DiscordAuthToken string

	recipients []MessageRecipient
	filters    []MessageFilter
}

// AddRecipients adds recipients to record messages to
func (s *BasicMessageSource) AddRecipients(recipients ...MessageRecipient) {
	s.recipients = append(s.recipients, recipients...)
}

func (s BasicMessageSource) sendMessageToAllRecipients(m *discordgo.Message) error {
	for _, recipient := range s.recipients {
		err := recipient.AddMessage(m)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddFilters adds filters for messages before they are sent to recipients
func (s *BasicMessageSource) AddFilters(filters ...MessageFilter) {
	s.filters = append(s.filters, filters...)
}

func (s BasicMessageSource) applyFilters(m *discordgo.Message, c *discordgo.Channel) bool {
	for _, filter := range s.filters {
		if !filter(m, c) {
			return false
		}
	}
	return true
}

// StreamMessages begins an async  stream of messages to the current stream of recipients
// Further changes to the recipient or filter list will result in undefined behavior
func (s BasicMessageSource) StreamMessages() <-chan Progress {
	pc := make(chan Progress)

	go s.asyncStreamMessages(pc)

	return pc
}

func (s BasicMessageSource) asyncStreamMessages(progressChan chan<- Progress) {
	session, err := discordgo.New(s.DiscordAuthToken)
	if err != nil {
		progressChan <- Progress{Error: err}
		return
	}

	channels, err := session.UserChannels()
	if err != nil {
		progressChan <- Progress{Error: err}
		return
	}

	var (
		messagesRead     = 0
		messagesRecorded = 0
		totalMessages    = 0
	)

	for _, ch := range channels {
		totalMessages += len(ch.Messages)
	}

	for _, ch := range channels {
		for _, msg := range ch.Messages {
			if s.applyFilters(msg, ch) {
				err = s.sendMessageToAllRecipients(msg)
				messagesRecorded++
			} else {
				err = nil
			}
			messagesRead++
			progressChan <- Progress{
				MessagesRead:     messagesRead,
				MessagesRecorded: messagesRecorded,
				PercentComplete:  float64(totalMessages) / float64(messagesRead),
				Error:            err,
			}
		}
	}

	close(progressChan)
}
