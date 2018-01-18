package discordstats

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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

func (s BasicMessageSource) sendMessageToAllRecipients(m *CrawledMessage) error {
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

func (s BasicMessageSource) applyFilters(m *CrawledMessage) bool {
	for _, filter := range s.filters {
		if !filter(m) {
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
		progressChan <- Progress{Error: fmt.Errorf("error creating discord session: %v", err)}
		close(progressChan)
		return
	}

	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		progressChan <- Progress{Error: fmt.Errorf("error getting discord guilds: %v", err)}
		close(progressChan)
		return
	}

	channels := []*discordgo.Channel{}

	for _, guild := range guilds {
		var guildChannels []*discordgo.Channel
		guildChannels, err = session.GuildChannels(guild.ID)
		if err != nil {
			progressChan <- Progress{Error: fmt.Errorf("error getting guild channels: %v", err)}
			close(progressChan)
			return
		}
		channels = append(channels, guildChannels...)
	}

	privateChannels, err := session.UserChannels()
	if err != nil {
		progressChan <- Progress{Error: fmt.Errorf("error getting discord user channels: %v", err)}
		close(progressChan)
		return
	}
	channels = append(channels, privateChannels...)

	var (
		messagesRead     = 0
		messagesRecorded = 0
		totalMessages    = 0
	)

	for _, ch := range channels {
		totalMessages += len(ch.Messages)
	}

	for _, ch := range channels {
		for crawlMsg := range crawlLinkedMessages(session, ch) {
			if messagesRecorded > 2000 {
				break
			}
			if crawlMsg.Error != nil {
				err = fmt.Errorf("error crawling messages: %v", err)
			} else if s.applyFilters(crawlMsg.LinkedMessage) {
				err = s.sendMessageToAllRecipients(crawlMsg.LinkedMessage)
				if err != nil {
					err = fmt.Errorf("error passing message to recipients: %v", err)
				}
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

type crawlMessage struct {
	LinkedMessage *CrawledMessage
	Error         error
}

func crawlLinkedMessages(s *discordgo.Session, channel *discordgo.Channel) <-chan crawlMessage {
	msgChan := make(chan crawlMessage)

	go func() {
		beforeID := ""
		var oldest *CrawledMessage
		for {
			messages, err := s.ChannelMessages(channel.ID, 100, beforeID, "", "")
			if err != nil {
				msgChan <- crawlMessage{Error: fmt.Errorf("error getting discord user channels: %v", err)}
				close(msgChan)
				return
			}
			linkedMessages := linkMessages(channel, messages, oldest)

			if oldest != nil {
				msgChan <- crawlMessage{LinkedMessage: oldest}
			}
			for i := 0; i < len(linkedMessages)-1; i++ {
				msgChan <- crawlMessage{LinkedMessage: linkedMessages[i]}
			}

			if len(messages) == 0 {
				break
			}
			oldest = linkedMessages[len(linkedMessages)-1]
			beforeID = oldest.Message.ID
		}

		if oldest != nil {
			msgChan <- crawlMessage{LinkedMessage: oldest}
		}
		close(msgChan)
	}()

	return msgChan
}

func linkMessages(channel *discordgo.Channel, messages []*discordgo.Message, capNewer *CrawledMessage) []*CrawledMessage {
	linkedMessages := []*CrawledMessage{}
	if len(messages) == 0 {
		return linkedMessages
	}
	for _, m := range messages {
		linkedMessages = append(linkedMessages, &CrawledMessage{Message: m, Channel: channel})
	}
	for i := 1; i < len(linkedMessages)-1; i++ {
		linkedMessages[i].Newer = linkedMessages[i-1]
		linkedMessages[i].Older = linkedMessages[i+1]
	}

	linkedMessages[0].Newer = capNewer
	if capNewer != nil {
		capNewer.Older = linkedMessages[0]
	}

	return linkedMessages
}
