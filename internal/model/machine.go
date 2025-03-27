package model

import (
	"encoding/json"
	"io"
)

const (
	requestMagic  ObjectName = "request"
	responseMagic ObjectName = "response"
)

// Response is the machine response.
type Response struct {
	Type    ObjectName `json:"type"`
	Message string     `json:"response"`
	IsError bool       `json:"error"`
}

// NewResponse builds a response.
func NewResponse(message string, isError bool) Response {
	return Response{
		Message: message,
		IsError: isError,
	}
}

// Encode is the response encoder.
func (r Response) Encode(writer io.Writer) error {
	r.Type = responseMagic

	return encode(writer, r)
}

// DecodeResponse is the response decoder.
func DecodeResponse(data string) *Response {
	var response Response
	if err := json.Unmarshal([]byte(data), &response); err == nil {
		if response.Type != responseMagic {
			return nil
		}

		return &response
	}

	return nil
}

// Request is an order passed to the machine.
type Request struct {
	Type    ObjectName `json:"type"`
	Message string     `json:"request"`
}

// NewRequest builds a request.
func NewRequest(message string) Request {
	return Request{
		Message: message,
	}
}

// Encode is the request encoder.
func (r Request) Encode(writer io.Writer) error {
	r.Type = requestMagic

	return encode(writer, r)
}

// DecodeRequest is the request decoder.
func DecodeRequest(data string) *Request {
	var request Request
	if err := json.Unmarshal([]byte(data), &request); err == nil {
		if request.Type != requestMagic {
			return nil
		}

		return &request
	}

	return nil
}
