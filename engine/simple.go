package engine

import (
	"crawler/fetcher"
	"crawler/model"
	"log"
)

type SimpleEngine struct{}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		parseResult, err := worker(r)
		if err != nil {
			continue
		}
		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			if profile, ok := item.(model.Profile); ok {
				log.Printf("user message %v \n", profile)
			} else {
				log.Printf("Got item %v \n", item)
			}
		}
	}
}

func worker(r Request) (ParserResult, error) {
	log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher:error fetching url:%s error:%v\n", r.Url, err)
		return ParserResult{}, nil
	}

	parseResult := r.ParserFunc(body)
	return parseResult, nil
}
