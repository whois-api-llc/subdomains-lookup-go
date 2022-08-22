package example

import (
	"context"
	"errors"
	subdomainslookup "github.com/whois-api-llc/subdomains-lookup-go"
	"log"
	"time"
)

func GetData(apikey string) {
	client := subdomainslookup.NewBasicClient(apikey)

	// Get parsed Subdomains Lookup API response as a model instance.
	subdomainsLookupResp, resp, err := client.Get(context.Background(), "whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only.
		subdomainslookup.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *subdomainslookup.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print some values from each returned record.
	for _, obj := range subdomainsLookupResp.Result.Records {
		log.Printf("Domain: %s, FirstSeen: %s, LastSeen: %s\n",
			obj.Domain,
			time.Unix(obj.FirstSeen, 0).Format(time.RFC3339),
			time.Unix(obj.LastSeen, 0).Format(time.RFC3339),
		)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := subdomainslookup.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com",
		subdomainslookup.OptionOutputFormat("JSON"))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
