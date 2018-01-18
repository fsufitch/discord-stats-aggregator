package main

import (
	"fmt"
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
	})

	msgCounter := pb.New(0)
	msgCounter.ShowBar = false
	msgCounter.ShowPercent = false
	msgCounter.ShowCounters = true
	msgCounter.ShowSpeed = true
	msgCounter.Start()

	for progress := range progressChan {
		msgCounter.Set(progress.MessagesRead)
	}
	msgCounter.Finish()

	result := <-resultChan
	fmt.Println("error was:", result.Error)
	fmt.Println(string(result.Data))
}
