package main

import (
	"context"

	"github.com/hashicorp/go-memdb"
)

var (
	ctx = context.Background()
)

var (
	schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"uuid": {
						Name:    "uuid",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "UUID"},
					},
					"uid": {
						Name:    "uid",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "UID"},
					},
					"uidNumber": {
						Name:    "uidNumber",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "UIDNumber"},
					},
					"gidNumber": {
						Name:    "gidNumber",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "GIDNumber"},
					},
				},
			},
		},
	}
)
