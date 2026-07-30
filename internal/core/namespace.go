package core

import (
	"context"
	"fmt"
)

func (c *Catalog) CreateNamespace(name string) (int64, error) {
	oid, err := c.q.CreateNamespace(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("create namespace %q: %w", name, err)
	}
	return oid, nil
}

func (c *Catalog) NamespaceOID(name string) (int64, error) {
	oid, err := c.q.NamespaceOID(context.Background(), name)
	if err != nil {
		return 0, fmt.Errorf("namespace %q: %w", name, err)
	}
	return oid, nil
}

type NamespaceInfo struct {
	OID  int64
	Name string
}

func (c *Catalog) Namespaces() ([]NamespaceInfo, error) {
	rows, err := c.q.ListNamespaces(context.Background())
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}
	out := make([]NamespaceInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, NamespaceInfo{OID: r.Oid, Name: r.Name})
	}
	return out, nil
}
