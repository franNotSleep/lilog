package domain

import (
	"testing"
	"time"
)

func TestNewInvoiceResponse(t *testing.T) {
	invResponse := NewInvoiceResponse(200, map[string]string{
		"content-type": "application/json; charset=utf-8",
	})

	if invResponse.StatusCode != 200 {
		t.Errorf("Unexpected Status Code. Expected=%d; Got=%d", 200, invResponse.StatusCode)
	}

	if val := invResponse.Headers["content-type"]; val != "application/json; charset=utf-8" {
		t.Errorf("Unexpected Content-type. Expected=%s; Got=%s", "application/json; charset=utf-8", val)
	}
}

func TestNewInvoiceRequest(t *testing.T) {
	invRequest := NewInvoiceRequest("POST", "ping/", map[string]string{}, map[string]string{}, map[string]string{}, "127.0.0.1", 2040)

	if invRequest.Method != "POST" {
		t.Errorf("Unexpected Method. Expected=%s; Got=%s", "POST", invRequest.Method)
	}

	if invRequest.URL != "ping/" {
		t.Errorf("Unexpected URL. Expected=%s; Got=%s", "ping/", invRequest.URL)
	}

	if invRequest.RemoteAddress != "127.0.0.1" {
		t.Errorf("Unexpected RemoteAddress. Expected=%s; Got=%s", "127.0.0.1", invRequest.RemoteAddress)
	}

	if invRequest.RemotePort != 2040 {
		t.Errorf("Unexpected RemotePort. Expected=%d; Got=%d", 2040, invRequest.RemotePort)
	}
}

func TestNewInvoice(t *testing.T) {
	inv := NewInvoice(time.Now().UnixMilli(), "info", 23231, "localhost", 5, "message", InvoiceRequest{}, InvoiceResponse{})

	if inv.Level != "info" {
		t.Errorf("Unexpected Level. Expected=%s; Got=%s", "info", inv.Level)
	}

	if inv.Hostname!= "localhost" {
		t.Errorf("Unexpected Hostname. Expected=%s; Got=%s", "hostname", inv.Hostname)
	}
}
