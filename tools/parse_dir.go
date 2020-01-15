package tools

import (
	"errors"
	"strings"
)

func ParseDir(dir, split string) (string, error) {

	var dl = make([]string, 0)

	if strings.Contains(dir, "/") {
		subDirs := strings.Split(dir, "/")

		for _, sd := range subDirs {
			if sd == "" {
				continue
			}

			if !strings.Contains(sd, "[") && !strings.Contains(sd, "]") {
				dl = append(dl, sd)
			} else {
				pre, suf := strings.Index(sd, "["), strings.Index(sd, "]")
				format := sd[pre+1 : suf]

				tf, err := GetLastTime(format, split)
				if err != nil {
					return "", err
				}
				dl = append(dl, tf)
			}
		}
	}

	if len(dl) == 0 {
		return "", errors.New("parse dir failed")
	}

	var s string
	for idx := 0; idx < len(dl); idx++ {
		s = s + "/" + dl[idx]
	}
	return s, nil
}
