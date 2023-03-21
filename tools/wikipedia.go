package tools

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var WikipediaTool = Tool{
	Name: "wikipedia",
	Description: `Search Wikipedia for the given term and return the first result
Usage: wikipedia <term>	`,
	Run: RunWikipedia,
}

// WikipediaResult is the result of a Wikipedia search
type WikipediaResult struct {
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

// RunWikipedia queries Wikipedia for the given search term and returns the
// snippet of the first result
func RunWikipedia(arg string) (string, error) {
	var result WikipediaResult
	queryParams := url.Values{
		"action":   []string{"query"},
		"list":     []string{"search"},
		"srsearch": []string{arg},
		"format":   []string{"json"},
		"srliimit": []string{"1"},
	}
	resp, err := http.Get("http://en.wikipedia.org/w/api.php?" + queryParams.Encode())
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
	snippet := result.Query.Search[0].Snippet
	snippet = strings.Replace(snippet, "<span class=\"searchmatch\">", "", -1)
	snippet = strings.Replace(snippet, "</span>", "", -1)
	return snippet, nil
}
