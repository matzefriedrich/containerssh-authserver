package resources

import (
	"embed"
)

//go:embed banner.txt postgres_schema.sql
var Resources embed.FS

const (
	BannerTxt         = "banner.txt"
	PostgresSchemaSql = "postgres_schema.sql"
)
