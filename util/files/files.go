package files

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func RelativePath(path1, path2 string) string {
	if strings.HasPrefix(path2, "/") || strings.HasPrefix(path2, "~/") {
		return path2
	}
	return filepath.ToSlash(filepath.Join(lang.If(strings.HasSuffix(path1, "/"), path1, filepath.Dir(path1)), path2))
}

func ResolveUserHomeDir(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.ToSlash(filepath.Join(optional.OfCommaErr(os.UserHomeDir()).Value(), strings.TrimPrefix(path, "~/")))
	}
	return path
}
