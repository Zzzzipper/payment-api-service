package third_party

import (
	"embed"
)

//go:embed OpenAPI/*
var OpenAPI embed.FS
var TestAPI embed.FS
