package tools

import (
	"fmt"
	"github.com/yitter/idgenerator-go/idgen"
)

func NextId(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, idgen.NextId())
}
