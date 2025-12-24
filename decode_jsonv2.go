//go:build goexperiment.jsonv2

package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
)

func (e *Event) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	panic("unimplemented")
}

var _ json.UnmarshalerFrom = (*Event)(nil)
