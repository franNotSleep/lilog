package server

import (
	"testing"
)

func TestMarshalUnMarshalReadRequest(t *testing.T) {
	rq := ReadReq{
		OpCode: OpRA,
		Server: "qvitae",
		From:   0,
		To:     0,
	}

	b, err := rq.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	rq2 := ReadReq{}
	err = rq2.UnmarshalBinary(b)

	if err != nil {
		t.Fatal(err)
	}

	if rq.OpCode != rq2.OpCode {
		t.Errorf("expected OpCode: %d; Got: %d\n", rq.OpCode, rq2.OpCode)
	}

	if rq.Server != rq2.Server {
		t.Errorf("expected Server: %s; Got: %s\n", rq.Server, rq2.Server)
	}

	if rq.From != rq2.From {
		t.Errorf("expected From: %d; Got: %d\n", rq.From, rq2.From)
	}

	if rq.To != rq2.To {
		t.Errorf("expected To: %d; Got: %d\n", rq.To, rq2.To)
	}
}
