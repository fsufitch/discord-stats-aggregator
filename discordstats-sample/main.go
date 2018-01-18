package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/fsufitch/discord-stats-aggregator"
	"github.com/fsufitch/discord-stats-aggregator/analyzer"
	"github.com/fsufitch/discord-stats-aggregator/filter"
)

func main() {
	authKey := os.Args[1]
	progressChan, resultChan := discordstats.EasyModeProgress(authKey, []discordstats.MessageFilter{
		filter.ByPublicChannel(true),
		filter.ByBot(false),
		filter.BeforeTime(time.Unix(1516221721, 0)),
	}, []discordstats.MessageAnalyzer{
		&analyzer.UserMessageTally{},
		&analyzer.ReactionTally{},
		&analyzer.TopReactionMessages{},
		&analyzer.TotalCount{},
		&analyzer.HereMentionTally{},
	})

	msgCounter := pb.New(0)
	msgCounter.ShowBar = false
	msgCounter.ShowPercent = false
	msgCounter.ShowCounters = true
	msgCounter.ShowSpeed = true
	msgCounter.Start()

	currentChannelName := ""

	for progress := range progressChan {
		msgCounter.Set(progress.MessagesRead)
		if progress.Error != nil {
			fmt.Fprintln(os.Stderr, "non-halting error:", progress.Error)
		}

		if progress.CurrentChannel.Name != currentChannelName {
			currentChannelName = progress.CurrentChannel.Name
			fmt.Fprintln(os.Stderr, "Now working on channel:", currentChannelName)
		}
	}
	msgCounter.Finish()

	result := <-resultChan
	if result.Error != nil {
		fmt.Fprintln(os.Stderr, "error was:", result.Error)
	} else {
		ioutil.WriteFile("output.json", result.Data, 0600)
	}
}
