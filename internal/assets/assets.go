package assets

import "embed"

//go:embed tongwen-dict/*.json
var Dicts embed.FS
