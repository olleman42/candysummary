package summarizer

// internal tally of total snack per consumer
type consumerEntry struct {
	name   string
	snacks map[string]int
}

// container for consumer-specific tally
type consumerMap map[string]*consumerEntry

func (cm consumerMap) appendEntry(entry HistoryEntry) {
	if consumer, ok := cm[entry.Name]; ok {
		if _, exists := consumer.snacks[entry.Candy]; exists {
			consumer.snacks[entry.Candy] += entry.Eaten
		} else {
			consumer.snacks[entry.Candy] = entry.Eaten
		}
	} else {
		cm[entry.Name] = &consumerEntry{name: entry.Name, snacks: map[string]int{entry.Candy: entry.Eaten}}
	}
}
