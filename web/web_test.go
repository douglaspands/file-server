package web_test

import (
	"io/fs"
	"testing"

	"github.com/douglas/file-server/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedWebFiles(t *testing.T) {
	t.Run("Given embedded templates FS When reading base layout Then file exists and is not empty", func(t *testing.T) {
		tmplFS := web.GetTemplatesFS()
		data, err := fs.ReadFile(tmplFS, "templates/layouts/base.html")

		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "File Server")
		assert.Contains(t, string(data), "icon.svg")
		assert.Contains(t, string(data), "theme-color")
		assert.Contains(t, string(data), "manifest.json")
	})

	t.Run("Given embedded templates FS When reading explorer page Then file exists", func(t *testing.T) {
		tmplFS := web.GetTemplatesFS()
		data, err := fs.ReadFile(tmplFS, "templates/pages/explorer.html")

		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "fileExplorer")
	})

	t.Run("Given embedded templates FS When reading gui launcher page Then file exists", func(t *testing.T) {
		tmplFS := web.GetTemplatesFS()
		data, err := fs.ReadFile(tmplFS, "templates/pages/gui_launcher.html")

		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "guiLauncher")
		assert.Contains(t, string(data), "icon.svg")
		assert.Contains(t, string(data), "theme-color")
		assert.Contains(t, string(data), "manifest.json")
	})

	t.Run("Given embedded static FS When reading static stylesheet Then file exists", func(t *testing.T) {
		staticFS, err := web.GetStaticFileSystem()
		require.NoError(t, err)

		file, err := staticFS.Open("/css/styles.css")
		require.NoError(t, err)
		defer file.Close()

		stat, err := file.Stat()
		require.NoError(t, err)
		assert.False(t, stat.IsDir())
	})

	t.Run("Given embedded static FS When reading app.js Then file exists", func(t *testing.T) {
		staticFS, err := web.GetStaticFileSystem()
		require.NoError(t, err)

		file, err := staticFS.Open("/js/app.js")
		require.NoError(t, err)
		defer file.Close()

		stat, err := file.Stat()
		require.NoError(t, err)
		assert.False(t, stat.IsDir())
	})

	t.Run("Given embedded static FS When reading manifest.json Then file exists and has theme_color", func(t *testing.T) {
		staticFS, err := web.GetStaticFileSystem()
		require.NoError(t, err)

		file, err := staticFS.Open("/manifest.json")
		require.NoError(t, err)
		defer file.Close()

		stat, err := file.Stat()
		require.NoError(t, err)
		assert.False(t, stat.IsDir())
	})

	t.Run("Given embedded static FS When reading icon assets Then files exist and are non-empty", func(t *testing.T) {
		staticFS, err := web.GetStaticFileSystem()
		require.NoError(t, err)

		svgFile, err := staticFS.Open("/assets/icon.svg")
		require.NoError(t, err)
		defer svgFile.Close()

		icoFile, err := staticFS.Open("/assets/icon.ico")
		require.NoError(t, err)
		defer icoFile.Close()
	})

	t.Run("Given embedded templates FS When reading file_table partial Then file exists and contains view and download links", func(t *testing.T) {
		tmplFS := web.GetTemplatesFS()
		data, err := fs.ReadFile(tmplFS, "templates/partials/file_table.html")

		require.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, string(data), "/view/")
		assert.Contains(t, string(data), "/download/")
		assert.Contains(t, string(data), "Visualizar no navegador")
		assert.Contains(t, string(data), "IsViewable")
	})
}
