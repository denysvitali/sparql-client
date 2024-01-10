package sparql

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Endpoint string
	client   *http.Client
	logger   *logrus.Logger
}

func (c *Client) Query(s string) (*Result, error) {
	data := url.Values{
		"query": []string{s},
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/sparql-results+json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code %d", resp.StatusCode)
	}

	// Parse JSON
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) SetLogger(logger *logrus.Logger) {
	if logger != nil {
		c.logger = logger
	}
}

// FetchAll fetches all the results of a query
// by incrementing the offset until there are no more results
func (c *Client) FetchAll(query string, limit int) (*Result, error) {
	res, err := c.fetchAll(query, 0, limit)
	return res, err
}

// fetchAll fetches all the results of a query
// this is the inner function that gets recursively called
func (c *Client) fetchAll(query string, offset int, limit int) (*Result, error) {
	c.logger.Debugf("fetching %d results from offset %d", limit, offset)
	finalQuery := fmt.Sprintf("%s\nOFFSET %d\nLIMIT %d", query, offset, limit)
	result, err := c.Query(finalQuery)
	if err != nil {
		return nil, err
	}

	finalResult := result

	if len(result.Results.Bindings) == limit {
		// There are more results, fetch them
		res, err := c.fetchAll(query, offset+limit, limit)
		if err != nil {
			return nil, err
		}
		finalResult.Results.Bindings = append(finalResult.Results.Bindings, res.Results.Bindings...)
	}

	return finalResult, nil
}

func New(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
		client:   http.DefaultClient,
		logger:   logrus.New(),
	}
}
