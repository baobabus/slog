// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Formatter func(string, []interface{}, []error) string

func SimpleFormatter(message string, v []interface{}, es []error) string {
	buf := &bytes.Buffer{}
	sep := ""
	if len(message) > 0 {
		buf.WriteString(message)
		sep = " "
	}
	for i := 0; i < len(v); i++ {
		buf.WriteString(sep)
		if i&1 == 0 {
			sep = "="
		} else {
			sep = " "
		}
		buf.WriteString(asString(v[i]))
	}
	if es != nil && len(es) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(" - ")
		}
		ml := len(es) > 1
		for _, e := range es {
			if ml {
				buf.WriteString("\n\t")
			}
			if e == errSuccess || e == errEllipsis {
				buf.WriteString(e.Error())
			} else {
				fmt.Fprintf(buf, "error=%v", e)
			}
		}
	}
	return buf.String()
}

func CompactJsonFormatter(message string, v []interface{}, e []error) string {
	return jsonFormatter(message, v, e, false)
}

func PrettyJsonFormatter(message string, v []interface{}, e []error) string {
	return jsonFormatter(message, v, e, true)
}

func jsonFormatter(message string, v []interface{}, e []error, pretty bool) string {
	m := make(map[string]interface{})
	for i := 0; i < len(v); i += 2 {
		k, ok := v[i].(string)
		if ok && i < (len(v)-1) {
			m[k] = v[i+1]
		}
	}
	if e != nil && len(e) > 0 {
		switch {
		case len(e) == 1 && (e[0] == nil || e[0] == errSuccess):
			m["success"] = true
		default:
			es := make([]string, 0, len(e))
			for _, v := range e {
				es = append(es, v.Error())
			}
			m["errors"] = es
		}
	}
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(v, "", "    ")
	} else {
		b, err = json.Marshal(m)
	}
	if err != nil {
		return message
	}
	if len(b) > 0 {
		if len(message) > 0 {
			return message + " " + string(b)
		} else {
			return string(b)
		}
	} else {
		return message
	}
}

func asString(value interface{}) string {
	if value == nil {
		return "<nil>"
	}
	switch v := value.(type) {
	case string:
		return v
	case time.Time:
		return v.Format(time.RFC3339)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%+v", value)
	}
}
