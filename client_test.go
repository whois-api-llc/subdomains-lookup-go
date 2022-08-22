package subdomainslookup

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathSubdomainsLookupResponseOK         = "/SubdomainsLookup/ok"
	pathSubdomainsLookupResponseError      = "/SubdomainsLookup/error"
	pathSubdomainsLookupResponse500        = "/SubdomainsLookup/500"
	pathSubdomainsLookupResponsePartial1   = "/SubdomainsLookup/partial"
	pathSubdomainsLookupResponsePartial2   = "/SubdomainsLookup/partial2"
	pathSubdomainsLookupResponseUnparsable = "/SubdomainsLookup/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Subdomains Lookup API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathSubdomainsLookupResponseOK:
		case pathSubdomainsLookupResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathSubdomainsLookupResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathSubdomainsLookupResponsePartial1:
			response = response[:len(response)-10]
		case pathSubdomainsLookupResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathSubdomainsLookupResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Subdomains Lookup API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:              apiServer.Client(),
		SubdomainsLookupBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestSubdomainsLookupGet tests the Get function.
func TestSubdomainsLookupGet(t *testing.T) {
	checkResultRec := func(res *SubdomainsLookupResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"search":"whoisxmlapi.com","result":{"count":4,"records":[
{"domain":"internet-retailers.whoisxmlapi.com","firstSeen":1585329674,"lastSeen":1635120000},
{"domain":"tools.whoisxmlapi.com","firstSeen":1548470536,"lastSeen":1656892800},
{"domain":"account-api.whoisxmlapi.com","firstSeen":1546617362,"lastSeen":1589075875},
{"domain":"drs-api.whoisxmlapi.com","firstSeen":1546628980,"lastSeen":1589089211}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathSubdomainsLookupResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathSubdomainsLookupResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathSubdomainsLookupResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathSubdomainsLookupResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathSubdomainsLookupResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathSubdomainsLookupResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("XML"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathSubdomainsLookupResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("SubdomainsLookup.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("SubdomainsLookup.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("SubdomainsLookup.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestSubdomainsLookupGetRaw tests the GetRaw function.
func TestSubdomainsLookupGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"search":"whoisxmlapi.com","result":{"count":4,"records":[
{"domain":"internet-retailers.whoisxmlapi.com","firstSeen":1585329674,"lastSeen":1635120000},
{"domain":"tools.whoisxmlapi.com","firstSeen":1548470536,"lastSeen":1656892800},
{"domain":"account-api.whoisxmlapi.com","firstSeen":1546617362,"lastSeen":1589075875},
{"domain":"drs-api.whoisxmlapi.com","firstSeen":1546628980,"lastSeen":1589089211}]}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathSubdomainsLookupResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathSubdomainsLookupResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathSubdomainsLookupResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathSubdomainsLookupResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathSubdomainsLookupResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("XML"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathSubdomainsLookupResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathSubdomainsLookupResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "invalid argument: \"domainName\" can not be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("SubdomainsLookup.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("SubdomainsLookup.Get() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
