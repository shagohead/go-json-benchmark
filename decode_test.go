//go:build goexperiment.jsonv2

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"
	"testing"

	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/go-faster/jx"
)

var eventData []byte

func init() {
	f, err := os.Open(path.Join("testdata", "event.json"))
	if err != nil {
		log.Fatal(err)
	}
	eventData, err = io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
}

func TestJXDecoder(t *testing.T) {
	d := jx.GetDecoder()
	d.ResetBytes(eventData)
	e := new(Event)
	if err := e.Decode(d); err != nil {
		t.Fatal(err)
	}
	t.Logf("jx decoded event: %+v", *e)
}

func BenchmarkJXDecoder(b *testing.B) {
	d := jx.GetDecoder()
	b.ReportAllocs()
	for b.Loop() {
		d.ResetBytes(eventData)
		if err := (&Event{}).Decode(d); err != nil {
			b.Fatal(err)
		}
	}
}

func TestJSONV2Decoder(t *testing.T) {
	d := jsontext.NewDecoder(bytes.NewReader(eventData))
	e := new(Event)
	if err := json.UnmarshalDecode(d, e); err != nil {
		t.Fatal(err)
	}
	t.Logf("jsonv2 decoded event: %+v", *e)
}

func BenchmarkJSONV2Decoder(b *testing.B) {
	r := bytes.NewReader(eventData)
	d := jsontext.NewDecoder(r)
	b.ReportAllocs()
	for b.Loop() {
		r.Reset(eventData)
		d.Reset(r)
		if err := json.UnmarshalDecode(d, (&Event{})); err != nil {
			b.Fatal(err)
		}
	}
}
