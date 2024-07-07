package options

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"

	"github.com/frannotsleep/lilog/types"
)

func TestServerNameDecode(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4852")
	if err != nil {
		t.Fatal(err)
	}

	server, err := net.ListenTCP("tcp", addr)

	if err != nil {
		t.Fatal(err)
	}

	defer server.Close()
	done := make(chan struct{})

	go func() {
		opts := LigLogOptions{}
		conn, err := server.AcceptTCP()

		if err != nil {
			t.Error(err)
		}

		defer conn.Close()

		var typ uint8
		err = binary.Read(conn, binary.BigEndian, &typ)

		if err != nil {
			t.Error(err)
		}

		if typ != ServerNameType {
			t.Errorf("invalid ServerNameType. got=%d", typ)
		}

		var payload types.Payload
		switch typ {
		case ServerNameType:
			payload = new(ServerName)
			opts.ServerName = payload.(*ServerName)
		default:
			t.Error("unknown type")
		}

		_, err = payload.ReadFrom(conn)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual("Web Api", opts.ServerName.String()) {
			t.Errorf("unexpected ServerName. got=%s\n", opts.ServerName)
		}

		t.Logf("%+v\n", opts)
		done <- struct{}{}
	}()

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		t.Fatal(err)
	}

	sn := ServerName("Web Api")
	var payload types.Payload = &sn

	_, err = payload.WriteTo(conn)

	if err != nil {
		t.Fatal(err)
	}
	<-done
}
