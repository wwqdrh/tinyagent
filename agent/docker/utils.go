package docker

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	// image
	ErrImageNotExist = errors.New("no this image")

	// container
	ErrInUsed       = errors.New("ErrInUsed")
	ErrObjectInUsed = errors.New("ErrObjectInUsed")
	ErrPortRules    = errors.New("ErrPortRules")

	// internal
	ErrCaptchaCode     = errors.New("ErrCaptchaCode")
	ErrAuth            = errors.New("ErrAuth")
	ErrRecordExist     = errors.New("ErrRecordExist")
	ErrRecordNotFound  = errors.New("ErrRecordNotFound")
	ErrStructTransform = errors.New("ErrStructTransform")
	ErrInitialPassword = errors.New("ErrInitialPassword")
	ErrNotSupportType  = errors.New("ErrNotSupportType")
	ErrInvalidParams   = errors.New("ErrInvalidParams")

	ErrTokenParse = errors.New("ErrTokenParse")
)

func stringsToMap(list []string) map[string]string {
	var lableMap = make(map[string]string)
	for _, label := range list {
		if strings.Contains(label, "=") {
			sps := strings.SplitN(label, "=", 2)
			lableMap[sps[0]] = sps[1]
		}
	}
	return lableMap
}
