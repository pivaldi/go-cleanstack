package migrations

import "embed"

// FS contains all embedded migration files (SQL and Go).
//
//go:embed *.sql *.go
var FS embed.FS
