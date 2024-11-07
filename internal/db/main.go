package yidb

import (
	"github.com/ostafen/clover"
	"path"
)

func CreatDB(name string) {
	dbPath := path.Join("./.db", name)
	clover.Open(dbPath)
}
