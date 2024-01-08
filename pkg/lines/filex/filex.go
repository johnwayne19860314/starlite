package filex

import (
	"os"
	"strings"

	"github.com/huandu/xstrings"

	"github.startlite.cn/itapp/startlite/pkg/lines/errorx"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errorx.WithStack(err)
	}
	return true, nil
}

func EnsureDir(dir string) {
	_ = os.MkdirAll(dir, 0777)
}

func GoFileName(name string) string {
	return strings.ToLower(xstrings.ToCamelCase(name))
}
