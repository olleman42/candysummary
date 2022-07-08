package summarizer

import "sort"

// Service - container type for summarizer's dependencies
type Service struct {
	client listEntryProvider
}

// New returns a new instance of the summarizer
func New(client listEntryProvider) *Service {
	return &Service{client}
}

// GetSummary returns an unserialized summary of all entries which are retrieved by the client
func (s Service) GetSummary() (Summary, error) {
	entries, err := s.client.GetEntries()
	if err != nil {
		return nil, err
	}

	cm := consumerMap{}
	// create a map of consumers with an reduced list of all consumed candy
	for _, entry := range entries {
		cm.appendEntry(entry)
	}

	// potential for optimization via fan-out, but overkill in this scenario
	summary := mapCmToSummary(cm)

	// sort the summary to get the top customers
	sort.SliceStable(summary, func(a, b int) bool {
		return summary[b].TotalSnacks < summary[a].TotalSnacks
	})

	return summary, nil
}

// SummaryEntry - the output type of the summary proccess
// this type could be tracked by a separate package too, but in this case it would be overkill as we hardly have to leave the context of the summary creation
type SummaryEntry struct {
	Name           string `json:"name"`
	FavouriteSnack string `json:"favouriteSnack"`
	TotalSnacks    int    `json:"totalSnacks"`
}

// Summary is an alias for a summary entry slice
type Summary []SummaryEntry

// HistoryEntry is a domain type to track each consumption instance
type HistoryEntry struct {
	Name  string
	Candy string
	Eaten int
}

type historyEntryContainer interface {
	Name() string
	Candy() string
	Eaten() int
}
type listEntryProvider interface {
	GetEntries() ([]HistoryEntry, error)
}

func mapCmToSummary(cm consumerMap) Summary {
	summary := make(Summary, len(cm))
	idx := 0

	for _, consumer := range cm {
		total := 0
		favSnack := ""
		for snack, count := range consumer.snacks {
			total += count
			if favSnack == "" {
				favSnack = snack
			}
			if count > consumer.snacks[favSnack] {
				favSnack = snack
			}
		}

		consumerSummary := SummaryEntry{
			Name:           consumer.name,
			FavouriteSnack: favSnack,
			TotalSnacks:    total,
		}

		summary[idx] = consumerSummary
		idx++
	}
	return summary
}
