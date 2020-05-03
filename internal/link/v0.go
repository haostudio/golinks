package link

import (
	"fmt"
	"strings"
)

// V0 returns a v0 link.
func V0(url string) Link {
	return Link{
		Version: 0,
		Blob:    []byte(strings.TrimRight(url, "/")),
	}
}

type v0 struct{}

func (v *v0) New(payload string) ([]byte, error) {
	return []byte(payload), nil
}

func (v *v0) Resolve(blob []byte, param string) (target string, err error) {
	target = string(blob)
	if len(param) > 0 {
		target = fmt.Sprintf("%s/%s", target, param)
	}
	return
}

func (v *v0) Describe(blob []byte) (string, error) {
	return string(blob), nil
}
