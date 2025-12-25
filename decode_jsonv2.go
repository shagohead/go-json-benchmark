//go:build goexperiment.jsonv2

package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"time"

	"github.com/go-faster/errors"
)

func (e *Event) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		key, err := d.ReadValue()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		switch string(key) {
		case `"sdk"`:
			err = e.SDK.UnmarshalJSONFrom(d)
		case `"platform"`:
			e.Platform, err = v2DecodeStr(d)
		case `"server_name"`:
			e.ServerName, err = v2DecodeStr(d)
		case `"environment"`:
			e.Environment, err = v2DecodeStr(d)
		case `"release"`:
			e.Release, err = v2DecodeStr(d)
		case `"level"`:
			e.Level, err = v2DecodeStr(d)
		case `"event_id"`:
			e.EventID, err = v2DecodeStr(d)
		case `"message"`:
			e.Message, err = v2DecodeStr(d)
		case `"contexts"`:
			err = v2DecodeMap(d, &e.Contexts, "")
		case `"extra"`:
			err = v2DecodeMap(d, &e.Extra, "")
		case `"user"`:
			err = v2DecodeMap(d, &e.User, "")
		case `"tags"`:
			err = v2DecodeStrMap(d, &e.Tags)
		case `"exception"`:
			err = e.Exception.UnmarshalJSONFrom(d)
		case `"timestamp"`:
			var tok jsontext.Token
			tok, err = d.ReadToken()
			if err != nil {
				break
			}
			e.Timestamp, err = time.Parse(time.RFC3339Nano, tok.String())
		default:
			err = d.SkipValue()
		}
		if err != nil {
			return errors.Wrap(err, string(key))
		}
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*Event)(nil)

func (s *SDK) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		key, err := d.ReadValue()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		switch string(key) {
		case `"name"`:
			s.Name, err = v2DecodeStr(d)
		case `"version"`:
			s.Version, err = v2DecodeStr(d)
		default:
			err = d.SkipValue()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*SDK)(nil)

func (f *Exceptions) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '[' {
		return fmt.Errorf("expected array, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != ']' {
		var exc Exception
		if err := (&exc).UnmarshalJSONFrom(d); err != nil {
			return errors.Wrapf(err, "[%d]", len(*f))
		}
		*f = append(*f, exc)
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*Exceptions)(nil)

func (e *Exception) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		key, err := d.ReadValue()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		switch string(key) {
		case `"module"`:
			e.Module, err = v2DecodeStr(d)
		case `"type"`:
			e.Type, err = v2DecodeStr(d)
		case `"value"`:
			e.Value, err = v2DecodeStr(d)
		case `"frames"`:
			err = e.Frames.UnmarshalJSONFrom(d)
		default:
			err = d.SkipValue()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*Exception)(nil)

func (f *Frames) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '[' {
		return fmt.Errorf("expected array, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != ']' {
		var frame Frame
		if err := (&frame).UnmarshalJSONFrom(d); err != nil {
			return errors.Wrapf(err, "[%d]", len(*f))
		}
		*f = append(*f, frame)
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*Frames)(nil)

func (f *Frame) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		key, err := d.ReadValue()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		switch string(key) {
		case `"filename"`:
			f.Filename, err = v2DecodeStr(d)
		case `"abs_path"`:
			f.AbsPath, err = v2DecodeStr(d)
		case `"module"`:
			f.Module, err = v2DecodeStr(d)
		case `"function"`:
			f.Function, err = v2DecodeStr(d)
		case `"lineno"`:
			f.LineNum, err = v2DecodeInt(d)
		case `"context_line"`:
			f.CtxLine, err = v2DecodeStr(d)
		case `"pre_context"`:
			err = v2DecodeContextLines(d, &f.PreCtx)
		case `"post_context"`:
			err = v2DecodeContextLines(d, &f.PreCtx)
		case `"in_app"`:
			f.InApp, err = v2DecodeBool(d)
		case `"vars"`:
			err = v2DecodeMap(d, &f.Vars, "")
		default:
			err = d.SkipValue()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
	}
	_, err := d.ReadToken()
	return err
}

var _ json.UnmarshalerFrom = (*Frame)(nil)

const emptyString = ""

func v2DecodeStr(d *jsontext.Decoder) (string, error) {
	tok, err := d.ReadToken()
	if err != nil {
		return emptyString, err
	}
	return tok.String(), nil
}

func v2DecodeInt(d *jsontext.Decoder) (int, error) {
	tok, err := d.ReadToken()
	if err != nil {
		return 0, err
	}
	return int(tok.Int()), nil // NOTE: On 32 systems value will be truncated
}

func v2DecodeBool(d *jsontext.Decoder) (bool, error) {
	tok, err := d.ReadToken()
	if err != nil {
		return false, err
	}
	return tok.Bool(), nil
}

func v2DecodeMap(d *jsontext.Decoder, dst *map[string]any, prefix string) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		if *dst == nil {
			*dst = make(map[string]any)
		}
		tok, err := d.ReadToken()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		key := tok.String()
		var v any
		switch d.PeekKind() {
		case '0':
			tok, err := d.ReadToken()
			if err != nil {
				break
			}
			v = tok.Int()
		case '{':
			err = v2DecodeMap(d, dst, fmt.Sprintf("%s%s.", prefix, key))
			if err == nil {
				continue
			}
		case '"':
			tok, err := d.ReadToken()
			if err != nil {
				break
			}
			v = tok.String()
		default:
			return d.SkipValue()
		}
		if err != nil {
			return errors.Wrap(err, key)
		}
		(*dst)[fmt.Sprintf("%s%s", prefix, key)] = v
	}
	_, err := d.ReadToken()
	return err
}

func v2DecodeStrMap(d *jsontext.Decoder, dst *map[string]string) error {
	if k := d.PeekKind(); k != '{' {
		return fmt.Errorf("expected object, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != '}' {
		if *dst == nil {
			*dst = make(map[string]string)
		}
		tok, err := d.ReadToken()
		if err != nil {
			return errors.Wrap(err, "reading object key")
		}
		key := tok.String()
		if k := d.PeekKind(); k != '"' {
			return fmt.Errorf("expected string, got %s", k)
		}
		tok, err = d.ReadToken()
		if err != nil {
			break
		}
		(*dst)[string(key)] = tok.String()
	}
	_, err := d.ReadToken()
	return err
}

func v2DecodeContextLines(d *jsontext.Decoder, dst *[]string) error {
	if k := d.PeekKind(); k != '[' {
		return fmt.Errorf("expected array, got %s", k)
	}
	if _, err := d.ReadToken(); err != nil {
		return err
	}
	for d.PeekKind() != ']' {
		tok, err := d.ReadToken()
		if err != nil {
			break
		}
		*dst = append(*dst, tok.String())
	}
	_, err := d.ReadToken()
	return err
}
