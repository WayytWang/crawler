package parser

import (
	"crawler/engine"
	"crawler/model"
	"regexp"
	"strconv"
)

var ageRe = regexp.MustCompile(`<div class="m-btn purple"[^>]+>([0-9]+)岁</div>`)
var heightRe = regexp.MustCompile(`<div class="m-btn purple"[^>]+>([0-9]+)cm</div>`)
var weightRe = regexp.MustCompile(`<div class="m-btn purple"[^>]+>([0-9]+)kg</div>`)
var incomeRe = regexp.MustCompile(`<div class="m-btn purple"[^>]+>月收入:([^<]+)</div>`)
var hukouRe = regexp.MustCompile(`<div class="m-btn pink"[^>]+>籍贯:([^<]+)</div>`)

func ParseProfile(contents []byte, name string) engine.ParserResult {
	profile := model.Profile{}

	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Age = age
	}

	height, err := strconv.Atoi(extractString(contents, heightRe))
	if err == nil {
		profile.Height = height
	}

	weight, err := strconv.Atoi(extractString(contents, weightRe))
	if err == nil {
		profile.Weight = weight
	}

	profile.Income = extractString(contents, incomeRe)
	profile.HuKou = extractString(contents, hukouRe)
	profile.Name = name

	result := engine.ParserResult{
		Items: []interface{}{profile},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
