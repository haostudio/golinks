package link

import (
	"fmt"
)

// implemented versions
var (
	versions = map[int]Version{
		0: &v0{},
		1: &v1{},
		2: &v2{},
	}
)

// Version defines the version interface, which is the operator for the blob.
type Version interface {
	New(payload string) (blob []byte, err error)
	Resolve([]byte, string) (string, error)
	Describe([]byte) (string, error)
}

// Link defines the link struct
type Link struct {
	Version int
	Blob    []byte
}

// New returns a new link of ver with payload.
func New(ver int, payload string) (*Link, error) {
	version, ok := versions[ver]
	if !ok {
		return nil, ErrVersionNotSupport
	}
	blob, err := version.New(payload)
	if err != nil {
		return nil, err
	}
	return &Link{
		Version: ver,
		Blob:    blob,
	}, nil
}

// GetRedirectLink returns the target redirect link of l.
func (l *Link) GetRedirectLink(param string) (string, error) {
	version, ok := versions[l.Version]
	if !ok {
		return "", ErrVersionNotSupport
	}
	return version.Resolve(l.Blob, param)
}

// Description returns the description string of l.
func (l *Link) Description() (string, error) {
	version, ok := versions[l.Version]
	if !ok {
		return "", ErrVersionNotSupport
	}
	desc, err := version.Describe(l.Blob)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("v%d|%s", l.Version, desc), nil
}
