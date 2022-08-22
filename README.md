[![subdomains-lookup-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![subdomains-lookup-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/subdomains-lookup-go)
[![subdomains-lookup-go test](https://github.com/whois-api-llc/subdomains-lookup-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/subdomains-lookup-go/actions/)

# Overview

The client library for
[Subdomains Lookup API](https://subdomains.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/subdomains-lookup-go
```

# Examples

Full API documentation available [here](https://subdomains.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := subdomainslookup.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := subdomainslookup.NewClient(apiKey, subdomainslookup.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Subdomains Lookup API lets you get a list of all subdomains related to a given domain name.

```go

// Make request to get all subdomains for the domain name.
subdomainsLookupResp, resp, err := client.Get(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

for _, obj := range subdomainsLookupResp.Result.Records {
    log.Printf("Domain: %s, FirstSeen: %s, LastSeen: %s\n", 
		obj.Domain, 
		time.Unix(obj.FirstSeen, 0).Format(time.RFC3339), 
		time.Unix(obj.LastSeen, 0).Format(time.RFC3339),
    )
}

// Make request to get raw data from Subdomains Lookup API.
resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
