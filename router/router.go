package router

import (
	"httpServer/converter"
	"os"
)

func HandleHttpRequest(payload converter.HttpPayload) converter.Response {
	response := converter.Response{HeaderMap: payload.HeaderMap, Protocol: payload.Protocol}
	response.HeaderMap["version"] = os.Getenv("VERSION")
	if payload.Path == "/healthz" {
		response.Status = 200
	} else {
		response.Status = 404
	}
	return response
}
