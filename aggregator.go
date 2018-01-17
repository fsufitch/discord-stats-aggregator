package discordstats

import "encoding/json"

// JSONDataAggregator combines data from multiple data providers and outputs JSON
type JSONDataAggregator struct {
	Providers map[string]DataProvider
}

// Aggregate gathers all the data from the providers and marshals it as JSON
func (a JSONDataAggregator) Aggregate() ([]byte, error) {
	dataContainer := map[string]interface{}{}
	for key, provider := range a.Providers {
		dataPiece, err := provider.GetData()
		if err != nil {
			return nil, err
		}
		dataContainer[key] = dataPiece
	}
	data, err := json.Marshal(dataContainer)
	return data, err
}
