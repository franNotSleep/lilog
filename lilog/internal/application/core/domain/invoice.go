package domain

type InvoiceRequest struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	Query         map[string]string `json:"query"`
	Params        map[string]string `json:"params"`
	Headers       map[string]string `json:"headers"`
	RemoteAddress string            `json:"remote_address"`
	RemotePort    int32             `json:"remote_port"`
}

type InvoiceResponse struct {
	StatusCode int32             `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

type Invoice struct {
	Time         int64    `json:"time"`
	Level        uint8   `json:"level"`
	PID          int32    `json:"pid"`
	Hostname     string   `json:"hostname"`
	InvoiceRequest      InvoiceRequest  `json:"request"`
	InvoiceResponse     InvoiceResponse `json:"response"`
	ResponseTime int32    `json:"response_time"`
	Message      string   `json:"message"`
}

func NewInvoiceResponse(statusCode int32, headers map[string]string) InvoiceResponse {
	return InvoiceResponse{StatusCode: statusCode, Headers: headers}
}

func NewInvoiceRequest(method string, url string, query map[string]string, params map[string]string, headers map[string]string, remoteAddress string, remotePort int32) InvoiceRequest {
	return InvoiceRequest{Method: method, URL: url, Query: query, Params: params, Headers: headers, RemoteAddress: remoteAddress, RemotePort: remotePort}
}

func NewInvoice(time int64, level uint8, pid int32, hostname string, responseTime int32, message string, request InvoiceRequest, response InvoiceResponse) Invoice {
	return Invoice{Time: time, Level: level, PID: pid, Hostname: hostname, InvoiceRequest: request, InvoiceResponse: response}
}
