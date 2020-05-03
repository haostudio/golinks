package link

import (
	"fmt"
	"strings"

	"github.com/haostudio/golinks/internal/encoding/gob"
)

// V2 returns a v2 link.
func V2(format string) (ln Link, err error) {
	payload := v2payload{
		VariableNum: parseV2NumVar(format),
		Format:      format,
	}
	blob, err := v2enc.Encode(payload)
	if err != nil {
		return
	}
	ln = Link{
		Version: 2,
		Blob:    blob,
	}
	return
}

type v2 struct{}

type v2payload struct {
	VariableNum int
	Format      string
}

func (v *v2) New(str string) ([]byte, error) {
	payload := v2payload{
		VariableNum: parseV2NumVar(str),
		Format:      str,
	}
	return v2enc.Encode(payload)
}

var v2enc = gob.New()

func (v *v2) Resolve(blob []byte, param string) (target string, err error) {
	var payload v2payload
	err = v2enc.Decode(blob, &payload)
	if err != nil {
		return
	}

	params := make([]string, 0, payload.VariableNum)
	for i := 0; i < payload.VariableNum; i++ {
		if len(param) == 0 {
			// Invalid number of parameters
			err = ErrInvalidParams
			return
		}
		var elem string
		elem, param = Pop(param, "/")
		params = append(params, elem)
	}

	target = payload.Format
	for i := 0; i < payload.VariableNum; i++ {
		target = strings.Replace(target, fmt.Sprintf("{%d}", i), params[i], -1)
	}
	return
}

func (v *v2) Describe(blob []byte) (string, error) {
	var payload v2payload
	err := v2enc.Decode(blob, &payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("var_num:%d|%s", payload.VariableNum, payload.Format), nil
}

func parseV2NumVar(url string) int {
	var i int
	for i < 1<<10 {
		if !strings.Contains(url, fmt.Sprintf("{%d}", i)) {
			return i
		}
		i++
	}
	return i
}
