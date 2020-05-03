package link

import (
	"strings"
)

// V1 returns a v1 link.
func V1(format string) Link {
	return Link{
		Version: 1,
		Blob:    []byte(format),
	}
}

type v1 struct{}

func (v *v1) New(payload string) ([]byte, error) {
	return []byte(payload), nil
}

func (v *v1) Resolve(blob []byte, param string) (target string, err error) {
	target = string(blob)
	target = strings.Replace(target, "{}", param, -1)
	return
}

func (v *v1) Describe(blob []byte) (string, error) {
	return string(blob), nil
}
