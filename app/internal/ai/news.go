package ai

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// NewsService fetches and caches military/USMC news headlines from Google News RSS.
type NewsService struct {
	mu      sync.Mutex
	cached  []NewsItem
	expires time.Time
}

// NewsItem is a single headline from the news feed.
type NewsItem struct {
	Title     string
	Source    string
	Published time.Time
	Link      string
}

// rssResponse maps the Google News RSS XML structure.
type rssResponse struct {
	Channel struct {
		Items []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Source  string `xml:"source"`
}

// Get returns cached news, fetching from Google News RSS if expired.
func (n *NewsService) Get() ([]NewsItem, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.cached != nil && time.Now().Before(n.expires) {
		return n.cached, nil
	}

	items, err := fetchNews()
	if err != nil {
		return nil, err
	}

	n.cached = items
	n.expires = time.Now().Add(30 * time.Minute)
	return items, nil
}

func fetchNews() ([]NewsItem, error) {
	// Multiple queries to get diverse military news relevant to a TBS XO
	queries := []string{
		"USMC OR \"Marine Corps\" OR \"Marines\"",
		"\"MCB Quantico\" OR \"Marine Corps Base Quantico\"",
		"military training OR \"Department of Defense\"",
	}

	seen := make(map[string]bool)
	var allItems []NewsItem

	client := &http.Client{Timeout: 10 * time.Second}

	for _, q := range queries {
		rssURL := "https://news.google.com/rss/search?q=" +
			url.QueryEscape(q) +
			"&hl=en-US&gl=US&ceid=US:en"

		resp, err := client.Get(rssURL)
		if err != nil {
			continue // try next query
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			continue
		}

		var rss rssResponse
		if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		for _, item := range rss.Channel.Items {
			// Google News titles often end with " - Source Name"
			title := item.Title
			source := item.Source
			if source == "" {
				if idx := strings.LastIndex(title, " - "); idx > 0 {
					source = title[idx+3:]
					title = title[:idx]
				}
			}

			if seen[title] {
				continue
			}
			seen[title] = true

			pubTime, _ := time.Parse(time.RFC1123, item.PubDate)
			if pubTime.IsZero() {
				pubTime, _ = time.Parse(time.RFC1123Z, item.PubDate)
			}

			allItems = append(allItems, NewsItem{
				Title:     title,
				Source:    source,
				Published: pubTime,
				Link:      item.Link,
			})
		}
	}

	if len(allItems) == 0 {
		return nil, fmt.Errorf("news fetch: no results from any query")
	}

	// Sort by recency (newest first), limit to 10
	for i := 0; i < len(allItems); i++ {
		for j := i + 1; j < len(allItems); j++ {
			if allItems[j].Published.After(allItems[i].Published) {
				allItems[i], allItems[j] = allItems[j], allItems[i]
			}
		}
	}

	if len(allItems) > 10 {
		allItems = allItems[:10]
	}

	return allItems, nil
}

// FormatNewsForPrompt formats news items for injection into a system prompt.
func FormatNewsForPrompt(items []NewsItem) string {
	if len(items) == 0 {
		return "No recent military news available."
	}

	var b strings.Builder
	b.WriteString("Recent Military/USMC News Headlines:\n")
	for i, item := range items {
		age := formatAge(item.Published)
		src := item.Source
		if src == "" {
			src = "News"
		}
		fmt.Fprintf(&b, "%d. %s (%s, %s)\n", i+1, item.Title, src, age)
	}
	return b.String()
}

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "recent"
	}
	d := time.Since(t)
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
