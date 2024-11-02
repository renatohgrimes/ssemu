package ssemu

import "embed"

//go:embed sql/*.sql
var Migrations embed.FS

//go:embed test/data.sql
var TestDataInsertCommand string

//go:embed web
var WebPages embed.FS
