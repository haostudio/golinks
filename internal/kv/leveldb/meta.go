package leveldb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/haostudio/golinks/internal/encoding/gob"
	"github.com/haostudio/golinks/internal/kv"
)

var (
	keyPrefix           = []byte("github.com/haostudio/golinks/internal/kv/leveldb.key")    // nolint: lll
	namespaceMetaPrefix = []byte("github.com/haostudio/golinks/internal/kv/leveldb.nsmeta") // nolint: lll

	metaEnc = gob.New()
)

type meta struct {
	// cache kv.Store
}

type namespaceMeta struct {
	Keys       map[string]struct{}
	Namespaces map[string]struct{}
}

func (m *meta) getKeyIn(key string, namespace ...string) []byte {
	k := append(keyPrefix[:0:0], keyPrefix...)
	k = append(k, []byte(strings.Join(append(namespace, key), "/"))...)
	return append(k, []byte("/"+key)...)
}

func (m *meta) addKeyIn(tx *leveldb.Transaction,
	key string, namespace ...string) ([]byte, error) {
	nsMeta, err := m.ensureNamespaceMeta(tx, nil, []string{}, namespace...)
	if err != nil {
		return nil, err
	}
	if nsMeta.Keys == nil {
		nsMeta.Keys = make(map[string]struct{})
	}
	nsMeta.Keys[key] = struct{}{}
	err = m.setNamespaceMeta(tx, nsMeta, namespace...)
	if err != nil {
		return nil, err
	}
	return m.getKeyIn(key, namespace...), nil
}

func (m *meta) deleteKeyIn(tx *leveldb.Transaction,
	key string, namespace ...string) ([]byte, error) {
	nsMeta, err := m.getNamespaceMetaTx(tx, namespace...)
	if errors.Is(err, kv.ErrNotFound) {
		return m.getKeyIn(key, namespace...), nil
	}
	if err != nil {
		return nil, err
	}
	if _, ok := nsMeta.Keys[key]; ok {
		delete(nsMeta.Keys, key)
		err = m.setNamespaceMeta(tx, nsMeta)
		if err != nil {
			return nil, err
		}
	}
	return m.getKeyIn(key, namespace...), nil
}

func (m *meta) getKeysIn(db *leveldb.DB, namespace ...string) (
	[]string, error) {
	nsMeta, err := m.getNamespaceMeta(db, namespace...)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(nsMeta.Keys))
	if len(keys) == 0 {
		return keys, nil
	}
	var i int
	for k := range nsMeta.Keys {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (m *meta) getKeysInTx(tx *leveldb.Transaction,
	namespace ...string) ([]string, error) {
	nsMeta, err := m.getNamespaceMetaTx(tx, namespace...)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(nsMeta.Keys))
	if len(keys) == 0 {
		return keys, nil
	}
	var i int
	for k := range nsMeta.Keys {
		keys[i] = k
		i++
	}
	return keys, nil
}

func (m *meta) getNamespacesIn(tx *leveldb.Transaction,
	namespace ...string) ([]string, error) {
	nsMeta, err := m.getNamespaceMetaTx(tx, namespace...)
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, len(nsMeta.Namespaces))
	if len(namespaces) == 0 {
		return namespaces, nil
	}
	var i int
	for ns := range nsMeta.Namespaces {
		namespaces[i] = ns
		i++
	}
	return namespaces, nil
}

func (m *meta) dropNamespaceMeta(tx *leveldb.Transaction,
	namespace ...string) error {
	// drop namespace meta
	k := m.getNamespaceMetaKey(namespace...)
	err := tx.Delete(k, &opt.WriteOptions{Sync: true})
	if err != nil {
		return err
	}

	// delete namespace in parent
	if len(namespace) == 0 {
		return m.setNamespaceMeta(tx, namespaceMeta{})
	}
	last := len(namespace) - 1
	nsMeta, err := m.getNamespaceMetaTx(tx, namespace[:last]...)
	if err != nil {
		return err
	}
	delete(nsMeta.Namespaces, namespace[last])
	return m.setNamespaceMeta(tx, nsMeta, namespace[:last]...)
}

func (m *meta) getNamespaceMeta(db *leveldb.DB,
	namespace ...string) (nsMeta namespaceMeta, err error) {
	k := m.getNamespaceMetaKey(namespace...)
	b, err := db.Get(k, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("%v: %w", err, kv.ErrNotFound)
		return
	}
	if err != nil {
		err = fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		return
	}
	err = metaEnc.Decode(b, &nsMeta)
	if err != nil {
		err = fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		return
	}
	return
}

func (m *meta) getNamespaceMetaTx(tx *leveldb.Transaction,
	namespace ...string) (nsMeta namespaceMeta, err error) {
	k := m.getNamespaceMetaKey(namespace...)
	b, err := tx.Get(k, nil)
	if errors.Is(err, leveldb.ErrNotFound) {
		err = fmt.Errorf("%v: %w", err, kv.ErrNotFound)
		return
	}
	if err != nil {
		err = fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		return
	}
	err = metaEnc.Decode(b, &nsMeta)
	if err != nil {
		err = fmt.Errorf("%v: %w", err, kv.ErrInternalError)
		return
	}
	return
}

func (m *meta) ensureNamespaceMeta(tx *leveldb.Transaction,
	root *namespaceMeta, rootNamespace []string, namespace ...string,
) (nsMeta namespaceMeta, err error) {
	var updateRoot bool
	if root == nil {
		var rootNsMeta namespaceMeta
		rootNsMeta, err = m.getNamespaceMetaTx(tx, rootNamespace...)
		if errors.Is(err, kv.ErrNotFound) {
			updateRoot = true
		} else if err != nil {
			return
		}
		root = &rootNsMeta
	}

	if len(namespace) == 0 {
		if updateRoot {
			err = m.setNamespaceMeta(tx, *root, rootNamespace...)
			if err != nil {
				return
			}
		}
		return *root, nil
	}
	if root.Namespaces == nil {
		root.Namespaces = make(map[string]struct{})
	}

	var nextNsMeta namespaceMeta
	subNamespace := append(rootNamespace[:0:0], rootNamespace...)
	subNamespace = append(subNamespace, namespace[0])

	_, found := root.Namespaces[namespace[0]]
	if found {
		nextNsMeta, err = m.getNamespaceMetaTx(tx, subNamespace...)
	} else {
		root.Namespaces[namespace[0]] = struct{}{}
		updateRoot = true
	}
	if err != nil {
		return
	}

	if updateRoot {
		err = m.setNamespaceMeta(tx, *root, rootNamespace...)
		if err != nil {
			return
		}
	}

	nextRoot := append(rootNamespace[:0:0], rootNamespace...)
	nextRoot = append(nextRoot, namespace[0])
	return m.ensureNamespaceMeta(tx, &nextNsMeta, nextRoot, namespace[1:]...)
}

func (m *meta) setNamespaceMeta(tx *leveldb.Transaction,
	nsMeta namespaceMeta, namespace ...string) error {
	b, err := metaEnc.Encode(&nsMeta)
	if err != nil {
		return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
	}

	k := m.getNamespaceMetaKey(namespace...)
	err = tx.Put(k, b, &opt.WriteOptions{Sync: true})
	if err != nil {
		return fmt.Errorf("%v: %w", err, kv.ErrInternalError)
	}
	return nil
}

func (m *meta) getNamespaceMetaKey(namespace ...string) []byte {
	k := append(namespaceMetaPrefix[:0:0], namespaceMetaPrefix...)
	return append(k, []byte(strings.Join(namespace, "/"))...)
}
