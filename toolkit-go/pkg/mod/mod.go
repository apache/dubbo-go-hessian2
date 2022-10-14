package mod

import (
	"encoding/json"
	"errors"
	"os/exec"
	"path"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/util"
)

type Version struct {
	Path    string
	Version string
}

type Mod struct {
	Module struct {
		Path string
	}
	Require []Version
	Replace []struct {
		Old Version
		New Version
	}
}

var (
	ErrNotMod = errors.New("not found go.mod")
)

func ReadMod(p string) (*Mod, error) {
	if len(p) == 0 {
		p = "./"
	}
	modFile := path.Join(p, "go.mod")
	if !util.FileExists(modFile) {
		return nil, ErrNotMod
	}
	// 读取数据
	cmd := exec.Command("go", "mod", "edit", "-json", modFile)
	d, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var dS Mod
	err = json.Unmarshal(d, &dS)
	if err != nil {
		return nil, err
	}
	return &dS, err
}
