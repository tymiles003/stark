// Copyright (C) 2014 Constantin Schomburg <me@cschomburg.com>
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proto

import (
	"reflect"
	"testing"
)

func TestMessageEncoding(t *testing.T) {
	m := Message{
		Version: VERSION,
		Id:      GenerateId(),
		Action:  "testaction",
		Source:  "testsource",
	}
	t.Log(m)

	enc, err := m.Encode()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(enc))

	dec, err := DecodeMessage(enc)
	t.Log(dec)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(dec, m) {
		t.Error("decoded message differs")
	}
}

func TestValid(t *testing.T) {
	var m Message
	m = Message{Version: VERSION, Id: "12345678", Action: "testaction", Source: "testsource"}
	if err := m.IsValid(); err != nil {
		t.Error(err)
	}

	m = Message{Version: VERSION, Id: "", Action: "testaction", Source: "testsource"}
	if err := m.IsValid(); err == nil {
		t.Error("Message without id passes as valid")
	}

	m = Message{Version: "", Id: "12345678", Action: "testaction", Source: "testsource"}
	if err := m.IsValid(); err == nil {
		t.Error("Message without version passes as valid")
	}

	m = Message{Version: VERSION, Id: "12345678", Action: "", Source: "testsource"}
	if err := m.IsValid(); err == nil {
		t.Error("Message without action passes as valid")
	}

	m = Message{Version: VERSION, Id: "12345678", Action: "testaction", Source: ""}
	if err := m.IsValid(); err == nil {
		t.Error("Message without source passes as valid")
	}
}

func TestReply(t *testing.T) {
	orig := Message{Version: VERSION, Id: GenerateId(), Action: "ping", Source: "originaldevice"}
	reply := orig.Reply(Message{Version: VERSION, Id: GenerateId(), Action: "ack", Source: "newdevice"})

	if reply.Id == orig.Id {
		t.Error("Reply has same id:", reply.Id)
	}
	if reply.CorrId != orig.Id {
		t.Error("Reply has wrong corrId:", reply.CorrId)
	}
	if reply.Destination != orig.Source {
		t.Error("Reply has wrong device:", reply.Destination)
	}
}

type simpleStruct struct {
	Key    string
	Number int
}

func (s simpleStruct) String() string {
	return "My key is " + s.Key
}

func TestEncodePayload(t *testing.T) {
	msg := Message{}
	exp := simpleStruct{"value", 35}
	if err := msg.EncodePayload(exp); err != nil {
		t.Fatal(err)
	}

	got := simpleStruct{}
	if err := msg.DecodePayload(&got); err != nil {
		t.Fatal(err)
	}

	if got.Key != "value" {
		t.Error("encode: wrong Key:", got.Key)
	}
	if msg.Text != "My key is value" {
		t.Error("encode: wrong text:", msg.Text)
	}

	if got != exp {
		t.Log(exp)
		t.Log(got)
		t.Error("decoded payload differs")
	}
}

func TestDecodeInvalidPayload(t *testing.T) {
	msg := Message{}
	var none struct{}
	if err := msg.DecodePayload(&none); err != nil {
		t.Fatal(err)
	}
}
