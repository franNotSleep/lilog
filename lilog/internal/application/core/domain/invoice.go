package domain

import "time"

type InvoiceRequest struct {
	Method        string            `json:"method" field:"method"`
	URL           string            `json:"url" field:"url"`
	Query         map[string]string `json:"query" field:"query"`
	Params        map[string]string `json:"params" field:"params"`
	Headers       map[string]string `json:"headers" field:"headers"`
	RemoteAddress string            `json:"remote_address" field:"remote_address"`
	RemotePort    int32             `json:"remote_port" field:"remote_port"`
}

type InvoiceResponse struct {
	StatusCode int32             `json:"status_code" field:"status_code"`
	Headers    map[string]string `json:"headers" field:"headers"`
}

type Invoice struct {
	ID              int64           `json:"id" field:"id"`
	Server          string          `json:"server" field:"server"`
	Time            time.Time       `json:"time" field:"time"`
	Level           uint8           `json:"level" field:"level"`
	PID             int32           `json:"pid" field:"pid"`
	Hostname        string          `json:"hostname" field:"hostname"`
	InvoiceRequest  InvoiceRequest  `json:"request" field:"request"`
	InvoiceResponse InvoiceResponse `json:"response" field:"response"`
	ResponseTime    int32           `json:"response_time" field:"response_time"`
	Message         string          `json:"message" field:"message"`
}

func NewInvoiceResponse(statusCode int32, headers map[string]string) InvoiceResponse {
	return InvoiceResponse{StatusCode: statusCode, Headers: headers}
}

func NewInvoiceRequest(method string, url string, query map[string]string, params map[string]string, headers map[string]string, remoteAddress string, remotePort int32) InvoiceRequest {
	return InvoiceRequest{Method: method, URL: url, Query: query, Params: params, Headers: headers, RemoteAddress: remoteAddress, RemotePort: remotePort}
}

func NewInvoice(t int64, level uint8, pid int32, hostname string, responseTime int32, message string, request InvoiceRequest, response InvoiceResponse) Invoice {
	return Invoice{Time: time.UnixMilli(t), Level: level, PID: pid, Hostname: hostname, Message: message, InvoiceRequest: request, InvoiceResponse: response, ResponseTime: responseTime}
}
