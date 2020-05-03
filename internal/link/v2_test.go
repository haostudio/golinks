package link

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseV2VarNum(t *testing.T) {
	cases := map[string]int{
		"zzzz.{0}":     1,
		"zzzz.{10}":    0,
		"zzzz.{1}.{0}": 2,
	}
	for url, num := range cases {
		require.Equal(t, num, parseV2NumVar(url))
	}
}
