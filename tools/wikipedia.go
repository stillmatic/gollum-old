package tools

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

var WikipediaTool = Tool{
	Name: "wikipedia",
	Description: `Search Wikipedia for the given term and return the first result
Usage: wikipedia <term>	`,
	Run: RunWikipedia,
}

// wikipediaResult is the result of a Wikipedia search
type wikipediaSearchResult struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Searchinfo struct {
			Totalhits         int    `json:"totalhits"`
			Suggestion        string `json:"suggestion"`
			Suggestionsnippet string `json:"suggestionsnippet"`
		} `json:"searchinfo"`
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    int       `json:"pageid"`
			Size      int       `json:"size"`
			Wordcount int       `json:"wordcount"`
			Snippet   string    `json:"snippet"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

type wikiExtractResult struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages []wikiPageResult `json:"pages"`
	} `json:"query"`
}

type wikiPageResult struct {
	Pageid  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

// RunWikipedia queries Wikipedia for the given search term and returns the
// summary of the first result
func RunWikipedia(arg string) (string, error) {
	// search for the term
	var result wikipediaSearchResult
	queryParams := url.Values{
		"action":   []string{"query"},
		"list":     []string{"search"},
		"srsearch": []string{arg},
		"format":   []string{"json"},
		"srliimit": []string{"1"},
	}
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?" + queryParams.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Query.Search) == 0 {
		return "", errors.New("no results")
	}
	// get the extract for the first result
	pageTitle := result.Query.Search[0].Title
	queryParams = url.Values{
		"action":        []string{"query"},
		"prop":          []string{"extracts"},
		"exsentences":   []string{"6"},
		"exlimit":       []string{"1"},
		"titles":        []string{pageTitle},
		"explaintext":   []string{"1"},
		"formatversion": []string{"2"},
		"format":        []string{"json"},
	}
	resp, err = http.Get("https://en.wikipedia.org/w/api.php?" + queryParams.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var extractResult wikiExtractResult
	if err := json.NewDecoder(resp.Body).Decode(&extractResult); err != nil {
		return "", err
	}

	// returns first page
	for _, page := range extractResult.Query.Pages {
		return page.Extract, nil
	}
	return "", errors.New("no results")
}
