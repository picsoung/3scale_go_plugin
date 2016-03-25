package go3scale

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	userAgentHeaderKey   = "X-3scale-User-Agent"
	userAgentHeaderValue = "plugin-golang-v#test"
)

// Client is identified by its ProviderKey and Host
type Client struct {
	ProviderKey, Host string
	httpClient        *http.Client
}

// NewClient Creates a new Client
func NewClient(providerKey, host string) (client *Client) {
	return &Client{
		ProviderKey: providerKey,
		Host:        host,
		httpClient:  http.DefaultClient,
	}
}

// Usage is a name and a value
type Usage struct {
	Name  string
	Value int
}

// AuthrepUserKey authenticates a request with userKey
func (client *Client) AuthrepUserKey(userKey string, usageArr []Usage) (Response, error) {
	req, err := http.NewRequest("GET", createURL(client, userKey, usageArr), nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Add(userAgentHeaderKey, userAgentHeaderValue)
	res, err := client.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		Code:    res.StatusCode,
		Message: string(body),
	}
	return response, nil
}

func createURL(client *Client, userKey string, usageArr []Usage) string {
	var buffer bytes.Buffer
	buffer.WriteString("https://")
	buffer.WriteString(client.Host)
	buffer.WriteString("/transactions/authrep.xml?user_key=")
	buffer.WriteString(userKey)
	buffer.WriteString("&provider_key=")
	buffer.WriteString(client.ProviderKey)
	for _, element := range usageArr {
		buffer.WriteString("&usage[")
		buffer.WriteString(element.Name)
		buffer.WriteString("]=")
		buffer.WriteString(strconv.Itoa(element.Value))
	}

	return buffer.String()
}

// Response object
type Response struct {
	Code    int
	Message string
}

// IsSuccess checks is response succeed
func (r *Response) IsSuccess() bool {
	return (r.Code == http.StatusOK || r.Code == http.StatusConflict)
}
