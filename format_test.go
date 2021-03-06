// Copyright 2016 Aleksey Blinov. All rights reserved.

package slog

import (
	"errors"
	"testing"
	"time"
)

type tcFormatter struct {
	msg string
	v   []interface{}
	err []error
	res string
}

func TestSimpleFormatter(tst *testing.T) {
	t1 := time.Date(2016, time.February, 21, 21, 3, 37, 0, time.UTC)
	t1s := "2016-02-21T21:03:37Z"
	for _, t := range []tcFormatter{
		{"", []interface{}{}, nil, ""},
		{"", []interface{}{}, []error{}, ""},
		{"", []interface{}{}, []error{errors.New("err")}, "error=err"},
		{"", []interface{}{}, []error{errSuccess}, "success"},
		{"", []interface{}{"foo"}, nil, "foo"},
		{"", []interface{}{"foo"}, []error{}, "foo"},
		{"", []interface{}{"foo"}, []error{errors.New("err")}, "foo - error=err"},
		{"", []interface{}{"foo"}, []error{errSuccess}, "foo - success"},
		{"", []interface{}{"foo", "bar"}, nil, "foo=bar"},
		{"", []interface{}{"foo", "bar"}, []error{}, "foo=bar"},
		{"", []interface{}{"foo", "bar"}, []error{errors.New("err")}, "foo=bar - error=err"},
		{"", []interface{}{"foo", "bar"}, []error{errSuccess}, "foo=bar - success"},
		{"msg", []interface{}{}, nil, "msg"},
		{"msg", []interface{}{}, []error{}, "msg"},
		{"msg", []interface{}{}, []error{errors.New("err")}, "msg - error=err"},
		{"msg", []interface{}{}, []error{errSuccess}, "msg - success"},
		{"msg", []interface{}{"foo"}, nil, "msg foo"},
		{"msg", []interface{}{"foo"}, []error{errors.New("err")}, "msg foo - error=err"},
		{"msg", []interface{}{"foo"}, []error{errSuccess}, "msg foo - success"},
		{"msg", []interface{}{"foo", "bar"}, nil, "msg foo=bar"},
		{"msg", []interface{}{"foo", "bar"}, []error{errors.New("err")}, "msg foo=bar - error=err"},
		{"msg", []interface{}{"foo", "bar"}, []error{errSuccess}, "msg foo=bar - success"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, nil, "msg foo=bar foo=bar"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, []error{errors.New("err")}, "msg foo=bar foo=bar - error=err"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, []error{errSuccess}, "msg foo=bar foo=bar - success"},
		{"msg", []interface{}{"v", nil}, nil, "msg v=<nil>"},
		{"msg", []interface{}{"v", errors.New("err")}, nil, "msg v=err"},
		{"msg", []interface{}{"v", 1}, nil, "msg v=1"},
		{"msg", []interface{}{"v", 0}, nil, "msg v=0"},
		{"msg", []interface{}{"v", -1}, nil, "msg v=-1"},
		{"msg", []interface{}{"v", 0.1}, nil, "msg v=0.1"},
		{"msg", []interface{}{"v", t1}, nil, "msg v=" + t1s},
		{"msg", []interface{}{"v", 10 * time.Second}, nil, "msg v=10s"},
	} {
		testFormatter(SimpleFormatter, &t, tst)
	}
}

func TestCompactJsonFormatter(tst *testing.T) {
	t1 := time.Date(2016, time.February, 21, 21, 3, 37, 0, time.UTC)
	t1s := "2016-02-21T21:03:37Z"
	for _, t := range []tcFormatter{
		{"", []interface{}{}, nil, "{}"},
		{"", []interface{}{}, []error{errors.New("err")}, "{\"errors\":[\"err\"]}"},
		{"", []interface{}{}, []error{errSuccess}, "{\"success\":true}"},
		{"", []interface{}{"foo"}, nil, "{}"},
		{"", []interface{}{"foo"}, []error{errors.New("err")}, "{\"errors\":[\"err\"]}"},
		{"", []interface{}{"foo"}, []error{errSuccess}, "{\"success\":true}"},
		{"", []interface{}{"foo", "bar"}, nil, "{\"foo\":\"bar\"}"},
		{"", []interface{}{"foo", "bar"}, []error{errors.New("err")}, "{\"errors\":[\"err\"],\"foo\":\"bar\"}"},
		{"", []interface{}{"foo", "bar"}, []error{errSuccess}, "{\"foo\":\"bar\",\"success\":true}"},
		{"msg", []interface{}{"v", nil}, nil, "msg {\"v\":null}"},
		{"msg", []interface{}{"v", 1}, nil, "msg {\"v\":1}"},
		{"msg", []interface{}{"v", 0}, nil, "msg {\"v\":0}"},
		{"msg", []interface{}{"v", -1}, nil, "msg {\"v\":-1}"},
		{"msg", []interface{}{"v", 0.1}, nil, "msg {\"v\":0.1}"},
		{"msg", []interface{}{"v", t1}, nil, "msg {\"v\":\"" + t1s + "\"}"},
		{"msg", []interface{}{"v", 10 * time.Second}, nil, "msg {\"v\":10000000000}"},
	} {
		testFormatter(CompactJsonFormatter, &t, tst)
	}
}

func testFormatter(fmtr Formatter, tc *tcFormatter, tst *testing.T) {
	res := fmtr(tc.msg, tc.v, tc.err)
	if tc.res != res {
		tst.Errorf("fail: expected \"%s\", but had \"%s\"", tc.res, res)
		return
	}
	if !testing.Short() {
		tst.Logf("pass: \"%s\"", res)
	}
}
