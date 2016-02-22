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
	err error
	res string
}

func TestSimpleFormatter(tst *testing.T) {
	t1 := time.Date(2016, time.February, 21, 21, 3, 37, 0, time.UTC)
	t1s := "2016-02-21T21:03:37Z"
	for _, t := range []tcFormatter{
		{"", []interface{}{}, nil, ""},
		{"", []interface{}{}, errors.New("err"), "error=err"},
		{"", []interface{}{}, errSuccess, "success"},
		{"", []interface{}{"foo"}, nil, "foo"},
		{"", []interface{}{"foo"}, errors.New("err"), "foo - error=err"},
		{"", []interface{}{"foo"}, errSuccess, "foo - success"},
		{"", []interface{}{"foo", "bar"}, nil, "foo=bar"},
		{"", []interface{}{"foo", "bar"}, errors.New("err"), "foo=bar - error=err"},
		{"", []interface{}{"foo", "bar"}, errSuccess, "foo=bar - success"},
		{"msg", []interface{}{}, nil, "msg"},
		{"msg", []interface{}{}, errors.New("err"), "msg - error=err"},
		{"msg", []interface{}{}, errSuccess, "msg - success"},
		{"msg", []interface{}{"foo"}, nil, "msg foo"},
		{"msg", []interface{}{"foo"}, errors.New("err"), "msg foo - error=err"},
		{"msg", []interface{}{"foo"}, errSuccess, "msg foo - success"},
		{"msg", []interface{}{"foo", "bar"}, nil, "msg foo=bar"},
		{"msg", []interface{}{"foo", "bar"}, errors.New("err"), "msg foo=bar - error=err"},
		{"msg", []interface{}{"foo", "bar"}, errSuccess, "msg foo=bar - success"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, nil, "msg foo=bar foo=bar"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, errors.New("err"), "msg foo=bar foo=bar - error=err"},
		{"msg", []interface{}{"foo", "bar", "foo", "bar"}, errSuccess, "msg foo=bar foo=bar - success"},
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
		{"", []interface{}{}, errors.New("err"), "{\"error\":\"err\"}"},
		{"", []interface{}{}, errSuccess, "{\"success\":true}"},
		{"", []interface{}{"foo"}, nil, "{}"},
		{"", []interface{}{"foo"}, errors.New("err"), "{\"error\":\"err\"}"},
		{"", []interface{}{"foo"}, errSuccess, "{\"success\":true}"},
		{"", []interface{}{"foo", "bar"}, nil, "{\"foo\":\"bar\"}"},
		{"", []interface{}{"foo", "bar"}, errors.New("err"), "{\"error\":\"err\",\"foo\":\"bar\"}"},
		{"", []interface{}{"foo", "bar"}, errSuccess, "{\"foo\":\"bar\",\"success\":true}"},
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
