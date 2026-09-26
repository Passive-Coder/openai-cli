package autocomplete

import (
	"slices"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestRebuildColonCommandBoundaries(t *testing.T) {
	t.Parallel()
	root := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "header", Aliases: []string{"H"}},
			&cli.BoolFlag{Name: "debug"},
		},
		Commands: []*cli.Command{
			{Name: "config:get", Aliases: []string{"cfg:get"}},
			{Name: "a::b"},
			{Name: "hidden:get", Hidden: true},
			{Name: "parent", Flags: []cli.Flag{&cli.StringFlag{Name: "value", Local: true}},
				Commands: []*cli.Command{{Name: "config:get"}}},
		},
	}
	for _, test := range []struct {
		name string
		args []string
		want []string
	}{
		{"nil", nil, nil},
		{"ordinary arguments", []string{"parent", "--debug"}, []string{"parent", "--debug"}},
		{"standalone colon", []string{"config", ":", "get"}, []string{"config:get"}},
		{"trailing colon", []string{"config:", "get"}, []string{"config:get"}},
		{"partial command", []string{"config:", "g"}, []string{"config:g"}},
		{"command alias", []string{"cfg:", "get"}, []string{"cfg:get"}},
		{"repeated standalone colons", []string{"a", ":", ":", "b"}, []string{"a::b"}},
		{"repeated trailing colons", []string{"a:", ":", "b"}, []string{"a::b"}},
		{"attached repeated colons", []string{"a::", "b"}, []string{"a::b"}},
		{"unknown command", []string{"unknown:", "get"}, []string{"unknown:", "get"}},
		{"hidden command", []string{"hidden:", "get"}, []string{"hidden:", "get"}},
		{"nested command", []string{"parent", "config:", "get"}, []string{"parent", "config:get"}},
		{"local flag value", []string{"parent", "--value", "config:", "get"}, []string{"parent", "--value", "config:", "get"}},
		{"inherited flag value", []string{"parent", "--header", "config:", "get"}, []string{"parent", "--header", "config:", "get"}},
		{"flag alias value", []string{"-H", "config:", "get"}, []string{"-H", "config:", "get"}},
		{"flag-like value", []string{"--header", "--debug", "config:", "get"}, []string{"--header", "--debug", "config:get"}},
		{"inline flag value", []string{"--header=config:", "get"}, []string{"--header=config:", "get"}},
		{"unknown value suffix", []string{"config:", "--debug"}, []string{"config:", "--debug"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := rebuildColonSeparatedArgs(root, test.args); !slices.Equal(got, test.want) {
				t.Fatalf("reconstructed %q as %q, want %q", test.args, got, test.want)
			}
		})
	}
}
