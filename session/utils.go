package session

import (
	"path"
	"strings"
)

func imgSrcToCharacter(src string) string {
	return strings.TrimSuffix(
		strings.TrimPrefix(path.Base(src), imgPathPrefix),
		path.Ext(src),
	)
}
