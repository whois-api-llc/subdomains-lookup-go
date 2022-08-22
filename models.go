package subdomainslookup

import (
	"fmt"
)

// Record contains info about subdomain: a domain name and timestamps.
type Record struct {
	// Domain is the domain name.
	Domain string `json:"domain"`

	// FirstSeen is the timestamp of the first time that the record was seen.
	FirstSeen int64 `json:"firstSeen"`

	// LastSeen is the timestamp of the last update for this record.
	LastSeen int64 `json:"lastSeen"`
}

// Result is the resulting object.
type Result struct {
	// Count is the number of records in result segment.
	Count int `json:"count"`

	// Records contains info about the resulting data: a domain name and timestamps.
	Records []Record `json:"records"`
}

// SubdomainsLookupResponse is a response of Subdomains Lookup API.
type SubdomainsLookupResponse struct {
	// Search is the target domain name.
	Search string `json:"search"`

	// Result is the resulting object.
	Result Result `json:"result"`
}

// ErrorMessage is an error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
