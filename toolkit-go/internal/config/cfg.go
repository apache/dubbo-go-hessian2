package config

import (
	"path"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/pkg/util"
)

// CustomStruct 用户自定义结构
type CustomStruct struct {
	JavaClassName string
	GoImportPath  string
	GoStructName  string
}

// HSFCfg hsf配置
type HSFCfg struct {
	CloseInternal bool
	CustomStruct  []CustomStruct
}
type MysqlTypeDefine struct {
	Type     string
	Unsigned []bool
	GoType   string
}

// MysqlCfg mysql 配置
type MysqlCfg struct {
	TypeConv []MysqlTypeDefine
}

// GotestsCfg go test配置
type GotestsCfg struct {
	MockSkipPkg []string
}

// Cfg 配置
type Cfg struct {
	Hsf     HSFCfg
	Gotests GotestsCfg
	Mysql   MysqlCfg
}

var (
	cfg     Cfg
	cfgPath string
	cfgDo   sync.Once
)

func SetConfigPath(p string) {
	cfgPath = p
}

func Config() Cfg {
	cfgDo.Do(func() {
		if len(cfgPath) == 0 {
			cfgPath = "./"
		}
		p := path.Join(cfgPath, ".gokit.toml")
		if util.FileExists(p) {
			_, err := toml.DecodeFile(p, &cfg)
			if err != nil {
				panic(".gokit.toml syntax error")
			}
		}
	})
	return cfg
}
