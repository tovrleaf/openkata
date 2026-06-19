package main

import "testing"

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no ansi",
			in:   "hello world",
			want: "hello world",
		},
		{
			name: "color codes",
			in:   "\x1b[32mgreen\x1b[0m text",
			want: "green text",
		},
		{
			name: "bold and reset",
			in:   "\x1b[1mbold\x1b[0m normal",
			want: "bold normal",
		},
		{
			name: "multiple sequences",
			in:   "\x1b[31m\x1b[1merror\x1b[0m: something failed",
			want: "error: something failed",
		},
		{
			name: "cursor movement",
			in:   "\x1b[2Kmoved\x1b[0m",
			want: "moved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripANSI(tt.in)
			if got != tt.want {
				t.Errorf("stripANSI(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCleanKiroOutput(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips header and ansi",
			in:   "> Using model claude-sonnet-4.6\n> Session started\n\x1b[32mHello\x1b[0m world",
			want: "Hello world",
		},
		{
			name: "preserves content lines",
			in:   "> Header\nline one\nline two\n",
			want: "line one\nline two",
		},
		{
			name: "empty input",
			in:   "",
			want: "",
		},
		{
			name: "only headers",
			in:   "> header1\n> header2\n",
			want: "",
		},
		{
			name: "mixed content",
			in:   "> info\nfirst\n> more info\nsecond",
			want: "first\nsecond",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanKiroOutput(tt.in)
			if got != tt.want {
				t.Errorf("cleanKiroOutput(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
