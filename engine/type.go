package engine

type Request struct {
	Url        string
	ParserFunc func([]byte) ParserResult
}

type ParserResult struct {
	Requests []Request
	Items    []interface{}
}

//NilParse 安全的nil，避免空指针异常
func NilParse(contents []byte) ParserResult {
	return ParserResult{}
}
