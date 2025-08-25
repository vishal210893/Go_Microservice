package repo

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit int `json:"limit" validate:"gte=1,lte=20"`
	Offset int `json:"offset" validate:"gte=0"`
	Sort string `json:"sort" validate:"oneof=asc desc"`
	Tags []string `json:"tags" validate:"max=5"`
	Search string `json:"search" validate:"max=100"`
	Since string `json:"since"`
	Until string `json:"until"`
}

// Parse extracts pagination parameters from HTTP request query string.
// It populates the PaginatedFeedQuery struct with values from the request URL parameters.
//
// Supported query parameters:
//   - limit: Number of items per page
//   - offset: Starting position for pagination
//   - sort: Sort order (asc/desc)
//   - tags: Comma-separated list of tags
//   - search: Search query string
//   - since: Start date filter (RFC3339 format)
//   - until: End date filter (RFC3339 format)
//
// Parameters:
//   - r: HTTP request containing query parameters
//
// Returns:
//   - PaginatedFeedQuery: Populated query struct with parsed parameters
//   - error: Parsing error if any parameter conversion fails
func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	// Parse limit parameter
	if limitStr := qs.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return fq, fmt.Errorf("invalid limit parameter: %w", err)
		}
		fq.Limit = limit
	}

	// Parse offset parameter
	if offsetStr := qs.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return fq, fmt.Errorf("invalid offset parameter: %w", err)
		}
		fq.Offset = offset
	}

	// Parse sort parameter
	if sortStr := qs.Get("sort"); sortStr != "" {
		fq.Sort = sortStr
	}

	// Parse tags parameter
	if tagsStr := qs.Get("tags"); tagsStr != "" {
		tags := strings.Split(tagsStr, ",")
		// Clean up tags (trim whitespace)
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		fq.Tags = tags
	} else {
		fq.Tags = []string{}
	}

	// Parse search parameter
	if searchStr := qs.Get("search"); searchStr != "" {
		fq.Search = strings.TrimSpace(searchStr)
	}

	// Parse since parameter
	if sinceStr := qs.Get("since"); sinceStr != "" {
		parsedTime, err := parseTime(sinceStr)
		if err != nil {
			return fq, fmt.Errorf("invalid since parameter: %w", err)
		}
		fq.Since = parsedTime
	}

	// Parse until parameter
	if untilStr := qs.Get("until"); untilStr != "" {
		parsedTime, err := parseTime(untilStr)
		if err != nil {
			return fq, fmt.Errorf("invalid until parameter: %w", err)
		}
		fq.Until = parsedTime
	}

	return fq, nil
}

// parseTime parses a time string and returns it in a standardized format.
// It accepts time strings in RFC3339 format and returns them normalized.
//
// Parameters:
//   - s: Time string in RFC3339 format (e.g., "2006-01-02 15:04:05")
//
// Returns:
//   - string: Normalized time string in RFC3339 format
//   - error: Parsing error if the time format is invalid
func parseTime(s string) (string, error) {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return "", fmt.Errorf("time must be in format '%s': %w", time.DateTime, err)
	}

	return t.Format(time.DateTime), nil
}