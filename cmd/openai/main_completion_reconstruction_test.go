package main

import (
	"strings"
	"testing"
)

func TestMainCompletionCommandBoundaries(t *testing.T) {
	for _, style := range []string{"bash", "zsh", "fish", "pwsh"} {
		t.Run(style, func(t *testing.T) {
			for _, test := range []struct {
				name string
				args []string
				want string
				code int
			}{
				{"unsplit command", []string{"chat:completions", "create", "--mo"}, "--model", 0},
				{"standalone colon", []string{"chat", ":", "completions", "create", "--mo"}, "--model", 0},
				{"trailing colon", []string{"chat:", "completions", "create", "--mo"}, "--model", 0},
				{"multiple trailing colons", []string{"admin:", "organization:", "users", "retrieve", "--user"}, "--user-id", 0},
				{"mixed colon boundaries", []string{"admin", ":", "organization:", "users", "retrieve", "--user"}, "--user-id", 0},
				{"partial command", []string{"chat:", "comp"}, "chat:completions", 0},
				{"empty command suffix", []string{"chat:", ""}, "chat:completions", 0},
				{"bool before command", []string{"--debug", "chat:", "completions", "create", "--mo"}, "--model", 0},
				{"value before command", []string{"--project", "synthetic-project", "chat:", "completions", "create", "--mo"}, "--model", 0},
				{"header before command", []string{"--header", "chat:", "completions", "create", "--mo"}, "--model", 0},
				{"header alias before command", []string{"-H", "chat:", "completions", "create", "--mo"}, "--model", 0},
				{"inline header before command", []string{"--header=chat:", "completions", "create", "--mo"}, "--model", 0},
				{"split header value", []string{"--header", "X", ":", "completions", "models", "retrieve", "--mo"}, "--model", 0},
				{"file before command", []string{"--mtls-client-key-file", "chat:", "completions", "create", "--mo"}, "--model", 0},
				{"value after command", []string{"beta:assistants", "create", "--instructions", "Prefix:", "--mo"}, "--model", 0},
				{"empty word after value", []string{"beta:assistants", "create", "--instructions", "Prefix:", ""}, "", 0},
				{"inherited value after command", []string{"chat:completions", "--header", "X:", "create", "--mo"}, "--model", 0},
				{"current header value", []string{"--header", "chat:"}, "", 11},
				{"current file value", []string{"--mtls-client-key-file", "chat:"}, "", 10},
			} {
				t.Run(test.name, func(t *testing.T) {
					args := append([]string{"openai", "__complete", "--"}, test.args...)
					got := runMainDispatch(t, style, args...)
					if got.code != test.code || got.stderr != "" {
						t.Fatalf("completion = %+v, want exit %d and no stderr", got, test.code)
					}
					want := test.want
					if style == "bash" && strings.Contains(want, ":") {
						want = want[strings.LastIndex(want, ":")+1:]
					} else if style == "zsh" {
						want = strings.ReplaceAll(want, ":", "\\:")
					}
					if want == "" {
						if got.stdout != "" {
							t.Fatalf("unexpected value suggestions: %q", got.stdout)
						}
						return
					}
					for _, line := range strings.Split(got.stdout, "\n") {
						if line == want || strings.HasPrefix(line, want+":") || strings.HasPrefix(line, want+"\t") {
							return
						}
					}
					t.Fatalf("completion = %q, want suggestion %q", got.stdout, want)
				})
			}
		})
	}
}
