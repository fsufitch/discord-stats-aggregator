package discordstats

// EasyMode wraps all discordstats functionality into one simple to use function (for the general case)
// Synchronous; only returns when finished, or errored
func EasyMode(authKey string, filters []MessageFilter, analyzers []MessageAnalyzer) ([]byte, error) {
	var source MessageSource = &BasicMessageSource{DiscordAuthToken: authKey}
	aggregator := JSONDataAggregator{}

	source.AddFilters(filters...)

	for _, analyzer := range analyzers {
		source.AddRecipients(analyzer)
		aggregator.Providers[analyzer.ID()] = analyzer
	}

	progress := source.StreamMessages()

	for p := range progress {
		if p.Error != nil {
			return nil, p.Error
		}
	}

	return aggregator.Aggregate()
}

// AsyncResult encapsulates the result of calling EasyModeProgress
type AsyncResult struct {
	Data  []byte
	Error error
}

// EasyModeProgress wraps all discordstats functionality into one simple to use function (for the general case)
// Asynchronous and provides progress updates (for a progress bar), and does not halt on errors
func EasyModeProgress(authKey string, filters []MessageFilter, analyzers []MessageAnalyzer) (<-chan Progress, <-chan AsyncResult) {
	var source MessageSource = &BasicMessageSource{DiscordAuthToken: authKey}
	aggregator := JSONDataAggregator{}

	source.AddFilters(filters...)

	for _, analyzer := range analyzers {
		source.AddRecipients(analyzer)
		aggregator.Providers[analyzer.ID()] = analyzer
	}

	progress := source.StreamMessages()
	ptProgress := make(chan Progress) // Passthrough
	resultChan := make(chan AsyncResult)

	go func() {
		for p := range progress {
			ptProgress <- p
		}

		data, err := aggregator.Aggregate()
		resultChan <- AsyncResult{data, err}
		close(resultChan)
	}()

	return ptProgress, resultChan
}
