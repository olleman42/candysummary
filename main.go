package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/olleman42/candysummary/client"
	"github.com/olleman42/candysummary/summarizer"
)

func main() {
	uri, ok := os.LookupEnv("CANDY_URI")
	if !ok {
		uri = "https://candystore.zimpler.net/"
	}
	summarizer := summarizer.New(client.New(uri))

	summary, err := summarizer.GetSummary()
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(summary)
}
