package sparql_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"

	"github.com/denysvitali/sparql-client"
)

func TestClientQuery(t *testing.T) {
	f, err := os.Open("./resources/test/mock-1.json")
	if err != nil {
		t.Fatal(err)
	}

	query := getText(t, "./resources/test/query-1.txt")

	gock.New("https://example.com").
		Post("/query").
		BodyString("query=.*").
		Reply(http.StatusOK).
		Body(f)

	client := sparql.New("https://example.com/query")
	res, err := client.Query(query)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Results.Bindings) != 50 {
		t.Fatalf("Expected 50 results, got %d", len(res.Results.Bindings))
	}
}

func TestClient_FetchAll(t *testing.T) {
	f, err := os.Open("./resources/test/mock-1.json")
	if err != nil {
		t.Fatal(err)
	}

	secondPage, err := os.Open("./resources/test/mock-2.json")
	if err != nil {
		t.Fatal(err)
	}

	query := getText(t, "./resources/test/query-2.txt")

	gock.New("https://example.com").
		Post("/query").
		BodyString("query=.*").
		Reply(http.StatusOK).
		Body(f)

	gock.New("https://example.com").
		Post("/query").
		BodyString("query=.*").
		Reply(http.StatusOK).
		Body(secondPage)

	client := sparql.New("https://example.com/query")
	res, err := client.FetchAll(query, 50)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Results.Bindings) != 50 {
		t.Fatalf("Expected 50 results, got %d", len(res.Results.Bindings))
	}
}

func getText(t *testing.T, s string) string {
	f, err := os.Open(s)
	if err != nil {
		t.Fatal(err)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
