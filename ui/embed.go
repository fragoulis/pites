package ui

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed all:build
var buildDir embed.FS

// DistDirFS contains the embedded dist directory files.
var (
	//nolint:gochecknoglobals
	DistDirFS = echo.MustSubFS(buildDir, "build")
)
