package main

import (
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

func (e *Event) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) (err error) {
		switch string(key) {
		case "sdk":
			err = e.SDK.Decode(d)
		case "platform":
			e.Platform, err = d.Str()
		case "server_name":
			e.ServerName, err = d.Str()
		case "environment":
			e.Environment, err = d.Str()
		case "release":
			e.Release, err = d.Str()
		case "level":
			e.Level, err = d.Str()
		case "event_id":
			e.EventID, err = d.Str()
		case "message":
			e.Message, err = d.Str()
		case "contexts":
			err = jxDecodeMap(d, &e.Contexts, "")
		case "extra":
			err = jxDecodeMap(d, &e.Extra, "")
		case "user":
			err = jxDecodeMap(d, &e.User, "")
		case "tags":
			err = jxDecodeStrMap(d, &e.Tags)
		case "exception":
			err = e.Exception.Decode(d)
		case "timestamp":
			var s string
			s, err = d.Str()
			if err != nil {
				break
			}
			e.Timestamp, err = time.Parse(time.RFC3339Nano, s)
		default:
			err = d.Skip()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
		return err
	})
}

func (s *SDK) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) (err error) {
		switch string(key) {
		case "name":
			s.Name, err = d.Str()
		case "version":
			s.Version, err = d.Str()
		default:
			err = d.Skip()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
		return err
	})
}

type Decoder interface {
	Decode(*jx.Decoder) error
}

// TODO: Use generic instead of concrete Exceptions and Frames.
type Array[T any, PT interface {
	*T
	Decoder
}] []T

func (a *Array[T, PT]) Decode(d *jx.Decoder) error {
	return d.Arr(func(d *jx.Decoder) error {
		var i T
		var pt PT = &i
		if err := pt.Decode(d); err != nil {
			return errors.Wrapf(err, "[%d]", len(*a))
		}
		*a = append(*a, i)
		return nil
	})
}

func (f *Exceptions) Decode(d *jx.Decoder) error {
	return d.Arr(func(d *jx.Decoder) error {
		var exc Exception
		if err := (&exc).Decode(d); err != nil {
			return errors.Wrapf(err, "[%d]", len(*f))
		}
		*f = append(*f, exc)
		return nil
	})
}

func (e *Exception) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) (err error) {
		switch string(key) {
		case "module":
			e.Module, err = d.Str()
		case "type":
			e.Type, err = d.Str()
		case "value":
			e.Value, err = d.Str()
		case "frames":
			err = e.Frames.Decode(d)
		default:
			err = d.Skip()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
		return err
	})
}

func (f *Frames) Decode(d *jx.Decoder) error {
	return d.Arr(func(d *jx.Decoder) error {
		var frame Frame
		if err := (&frame).Decode(d); err != nil {
			return errors.Wrapf(err, "[%d]", len(*f))
		}
		*f = append(*f, frame)
		return nil
	})
}

func (f *Frame) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) (err error) {
		switch string(key) {
		case "filename":
			f.Filename, err = d.Str()
		case "abs_path":
			f.AbsPath, err = d.Str()
		case "module":
			f.Module, err = d.Str()
		case "function":
			f.Function, err = d.Str()
		case "lineno":
			f.LineNum, err = d.Int()
		case "context_line":
			f.CtxLine, err = d.Str()
		case "pre_context":
			err = jxDecodeContextLines(d, &f.PreCtx)
		case "post_context":
			err = jxDecodeContextLines(d, &f.PreCtx)
		case "in_app":
			f.InApp, err = d.Bool()
		case "vars":
			err = jxDecodeMap(d, &f.Vars, "")
		default:
			err = d.Skip()
		}
		if err != nil {
			err = errors.Wrap(err, string(key))
		}
		return err
	})
}

func jxDecodeMap(d *jx.Decoder, dst *map[string]any, prefix string) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) (err error) {
		if *dst == nil {
			*dst = make(map[string]any)
		}
		var v any
		switch d.Next() {
		// case jx.Array:
		// case jx.Bool:
		// case jx.Invalid:
		// case jx.Null:
		case jx.Number:
			v, err = d.Int()
			if err != nil {
				v, err = d.Float64()
			}
		case jx.Object:
			err = jxDecodeMap(d, dst, fmt.Sprintf("%s%s.", prefix, key))
			if err == nil {
				return nil
			}
		case jx.String:
			v, err = d.Str()
		default:
			err = fmt.Errorf("unexpected jx.Type %s", d.Next().String())
		}
		if err != nil {
			return errors.Wrap(err, string(key))
		}
		(*dst)[fmt.Sprintf("%s%s", prefix, key)] = v
		return nil
	})
}

func jxDecodeStrMap(d *jx.Decoder, dst *map[string]string) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		if *dst == nil {
			*dst = make(map[string]string)
		}
		s, err := d.Str()
		if err != nil {
			return errors.Wrap(err, string(key))
		}
		(*dst)[string(key)] = s
		return err
	})
}

func jxDecodeContextLines(d *jx.Decoder, dst *[]string) error {
	return d.Arr(func(d *jx.Decoder) error {
		s, err := d.Str()
		if err != nil {
			return err
		}
		*dst = append(*dst, s)
		return nil
	})
}
