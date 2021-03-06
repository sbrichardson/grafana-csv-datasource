package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestParseCSV(t *testing.T) {
	for _, tt := range []struct {
		query  queryModel
		input  string
		output []*data.Field
	}{
		{query: queryModel{}, input: "foo,bar,baz\n1,2,3", output: []*data.Field{
			data.NewField("foo", nil, []string{"1"}),
		}},
	} {
		t.Run("", func(t *testing.T) {
			fields, err := parseCSV(tt.query, strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}

			f1 := data.Frame{Fields: fields}
			f2 := data.Frame{Fields: tt.output}

			b1, _ := f1.MarshalArrow()
			b2, _ := f2.MarshalArrow()

			if bytes.Equal(b1, b2) {
				t.Fatal("unexpected output")
			}
		})
	}
}

func TestParseTimeNaive(t *testing.T) {
	for _, tt := range []struct {
		input  string
		output string
	}{
		{input: "2018-08-19", output: "2006-01-02"},
		{input: "2018-08-19 12:11", output: "2006-01-02 15:04"},
		{input: "2018-08-19 12:11:35", output: "2006-01-02 15:04:05.999999"},
		{input: "2018-08-19 12:11:35.22", output: "2006-01-02 15:04:05.999999"},
		{input: "2018/08/19", output: "2006/01/02"},
		{input: "2018/08/19 12:11", output: "2006/01/02 15:04"},
		{input: "2018/08/19 12:11:35", output: "2006/01/02 15:04:05.999999"},
		{input: "2018/08/19 12:11:35.22", output: "2006/01/02 15:04:05.999999"},
		{input: "08/19/2018", output: "01/02/2006"},
		{input: "2018-07-05 12:54:00 UTC", output: "2006-01-02 15:04:05 MST"},
		{input: "2018-08-19 07:11:35.220 -05:00", output: "2006-01-02 15:04:05.999999 -07:00"},
		{input: "2018-08-19T12:11:35.220Z", output: "2006-01-02T15:04:05.999999Z"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			got, err := detectTimeLayoutNaive(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.output {
				t.Fatalf("want = %q; got = %q", tt.output, got)
			}
		})
	}
}
