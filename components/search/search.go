// Package search provides server-side helpers for search/filter components
// (CommandPalette, Combobox, SearchInput). Not a component package — no .templ files.
package search

import "strings"

// Item is a generic search result item.
type Item struct {
	Title       string
	Description string
	Path        string
	Section     string
	Keywords    []string
}

// Filter filters items by a search query, matching against title, description,
// section, and keywords. Results are sorted by relevance: exact title match
// first, then title prefix, then title contains, then keyword/description match.
// Limits results to maxItems (0 = unlimited).
func Filter(items []Item, query string, maxItems int) []Item {
	if query == "" {
		if maxItems > 0 && len(items) > maxItems {
			return items[:maxItems]
		}
		return items
	}

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return items
	}

	type scored struct {
		item  Item
		score int
	}
	var results []scored

	for _, item := range items {
		title := strings.ToLower(item.Title)
		desc := strings.ToLower(item.Description)
		section := strings.ToLower(item.Section)
		path := strings.ToLower(item.Path)

		if !strings.Contains(title, q) && !strings.Contains(desc, q) &&
			!strings.Contains(section, q) && !strings.Contains(path, q) &&
			!matchKeywords(item.Keywords, q) {
			continue
		}

		score := 0
		if title == q {
			score = 100
		} else if strings.HasPrefix(title, q) {
			score = 80
		} else if strings.Contains(title, q) {
			score = 60
		} else if strings.Contains(path, q) {
			score = 40
		} else {
			score = 20
		}
		if strings.Contains(desc, q) {
			score += 5
		}
		if matchKeywords(item.Keywords, q) {
			score += 10
		}
		results = append(results, scored{item, score})
	}

	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].score > results[i].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	result := make([]Item, len(results))
	for i, s := range results {
		result[i] = s.item
	}
	if maxItems > 0 && len(result) > maxItems {
		result = result[:maxItems]
	}
	return result
}

func matchKeywords(keywords []string, q string) bool {
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(kw), q) {
			return true
		}
	}
	return false
}

// HasMatches returns true if the query matches any item.
func HasMatches(items []Item, query string) bool {
	return len(Filter(items, query, 1)) > 0
}
