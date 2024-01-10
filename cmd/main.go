package main

import (
	"encoding/json"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
	"os"
	sparql "sparql-client"
)

var args struct {
	Endpoint string `arg:"-e,--endpoint,env:SPARQL_ENDPOINT,required"`
	Debug    *bool  `arg:"-D,--debug,env:SPARQL_CLIENT_DEBUG"`
	Query    string `arg:"-q,--query,env:SPARQL_QUERY,required"`
	FetchAll *bool  `arg:"-a,--fetch-all,env:SPARQL_FETCH_ALL"`
	Limit    int    `arg:"-l,--limit,env:SPARQL_LIMIT"`
	Output   string `arg:"-o,--output,env:SPARQL_OUTPUT"`
}

var logger = logrus.New()

func main() {
	arg.MustParse(&args)
	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	s := sparql.New(args.Endpoint)
	s.SetLogger(logger)

	var res *sparql.Result
	var err error

	if args.Limit <= 0 {
		args.Limit = 100
	}

	if args.FetchAll != nil && *args.FetchAll {
		res, err = s.FetchAll(
			readFromFile(args.Query),
			args.Limit,
		)
	} else {
		res, err = s.Query(
			readFromFile(args.Query),
		)
	}

	if err != nil {
		logger.Fatalf("unable to query: %v", err)
	}

	var enc *json.Encoder

	if args.Output == "" {
		enc = json.NewEncoder(os.Stdout)
	} else {
		f, err := os.Create(args.Output)
		if err != nil {
			logger.Fatalf("unable to create file: %v", err)
		}
		defer f.Close()
		enc = json.NewEncoder(f)
	}

	err = enc.Encode(res)
	if err != nil {
		logger.Fatalf("unable to encode result: %v", err)
	}
}
