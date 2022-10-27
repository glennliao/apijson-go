package service_test

import (
	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"path/filepath"
)

func init() {
	foundConfig()
}

type driver struct {
	*mysql.Driver
}

func foundConfig() {

	dir := "./"
	fileAdapter := g.Cfg().GetAdapter().(*gcfg.AdapterFile)

	for i := 0; i < 6; i++ {
		if gfile.Exists(filepath.Join(dir, "go.mod")) {
			fileAdapter.AddPath(gfile.Abs(dir))
			if gfile.Exists(filepath.Join(dir, "main.go")) {
				break
			}
		}

		dir = filepath.Join(dir, "../")
	}
}

func Pointer[T any](value T) *T {
	return &value
}
