package conf

import (
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/config"
)

func CheckInPkg(pkg string) bool {
	mockSkipPkg := config.Config().Gotests.MockSkipPkg
	if len(mockSkipPkg) == 0 {
		return false
	}
	for _, p := range mockSkipPkg {
		if p == pkg {
			return true
		}
	}
	return false
}
