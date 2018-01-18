package main

import (
	"fmt"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/fsufitch/discord-stats-aggregator"
	"github.com/fsufitch/discord-stats-aggregator/analyzers"
)

func main() {
	authKey := os.Args[1]
	progressChan, resultChan := discordstats.EasyModeProgress(authKey, []discordstats.MessageFilter{
		discordstats.FilterByPublicChannel(true),
		discordstats.FilterByBot(false),
	}, []discordstats.MessageAnalyzer{
		&analyzers.UserMessageCountAnalyzer{},
		&analyzers.ReactionCountAnalyzer{},
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
