package converter

import (
	"bytes"
	"strconv"
	"strings"
)

const (
	METHOD = "METHOD"
	PATH = "PATH"
	HEADER = "HEADER"
	BODY = "BODY"
	PROTOCOL = "PROTOCOL"
)

type ParseContext struct {

	HttpPayload

	sb strings.Builder

	headers []Header

	bodyLengh int

	Status string
}

type HttpPayload struct {

	Method string

	Path string

	HeaderMap map[string]string

	Ip string

	Body string

	Protocol string

}

type Header struct {
	key string
	value string
}

func (parseContext *ParseContext) Parse(part []byte) {

	if parseContext.Status == "" {
		parseContext.Status = METHOD
	}
	parseChar :=func (char rune) {
		switch parseContext.Status {
		case METHOD:
			parseContext.Status = parseContext.parseMethod(char)
		case PATH:
			parseContext.Status = parseContext.parsePath(char)
		case HEADER:
			parseContext.Status = parseContext.parseHeader(char)
		case BODY:
			parseContext.Status = parseContext.parseBody(char)
		case PROTOCOL:
			parseContext.Status = parseContext.parseProtocol(char)

		}
	}

	for _, char := range bytes.Runes(part) {
		parseChar(char)
	}



}
func (parseContext *ParseContext) parseBody(char rune) string {
	parseContext.bodyLengh ++
	parseContext.sb.WriteRune(char)
	length, _ := strconv.Atoi(parseContext.HeaderMap["Content-Length"])
	if parseContext.bodyLengh == length {
		parseContext.Body = parseContext.sb.String()
	}
	return BODY
}

func (parseContext *ParseContext) parseHeader(char rune) string {
	if parseContext.sb.Len() == 0 && char == '\r' {
		for _, v := range parseContext.headers {
			parseContext.HeaderMap[v.key] = v.value
		}
		return BODY
	} else if char == '\r' {
		headerParts := strings.Split(parseContext.sb.String(), ":")
		parseContext.headers = append(parseContext.headers, Header{key: strings.TrimSpace(headerParts[0]), value: strings.TrimSpace(headerParts[1])})
		parseContext.sb.Reset()
	} else if char != '\n' {
		parseContext.sb.WriteRune(char)
	}
	return HEADER
}

func (parseContext *ParseContext) parsePath(char rune) string {
	if char != ' ' {
		parseContext.sb.WriteRune(char)
		return PATH
	} else {
		parseContext.Path = parseContext.sb.String()
		parseContext.sb.Reset()
		return PROTOCOL
	}
}

func (parseContext *ParseContext) parseProtocol(char rune) string {
	if char != '\r' {
		parseContext.sb.WriteRune(char)
		return PROTOCOL
	} else {
		parseContext.Protocol = parseContext.sb.String()
		parseContext.sb.Reset()
		return HEADER
	}
}

func (parseContext *ParseContext) parseMethod(char rune) string {
	if char != ' ' {
		parseContext.sb.WriteRune(char)
		return METHOD
	} else {
		parseContext.Method = parseContext.sb.String()
		parseContext.sb.Reset()
		return PATH
	}
}