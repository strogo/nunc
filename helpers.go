package nunc

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func ResolvePath(file string) (path string) {
	return filepath.Join(is.Home, file)
}

func ResolveTaskID(id string) (context Context, taskId int64, err error) {
	clean := strings.Map(dropSpace, id)
	data := strings.SplitN(clean, "-", 2)
	switch len(data) {
	case 1:
		if strings.HasPrefix(data[0], "@") {
			err = fmt.Errorf("only context provided")
			return
		}
		data = append(data, data[0])
		data[0] = ""
		fallthrough
	case 2:
		if data[0] != "" {
			context, _, err = GetContext(data[0], true)
			if err != nil {
				return
			}
		}
		taskId, err = strconv.ParseInt(data[1], 0, 64)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("invalid task id")
		return
	}
	return
}

func dropSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}
