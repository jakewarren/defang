package defang

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Refang takes a defanged IOC and returns it to it's original form
func Refang(input interface{}) (string, error) {

	var output string

	switch i := input.(type) {
	case string:
		output = i
	case net.IP:
		output = i.String()
	case url.URL:
		output = i.String()
	default:
		return "", errors.New("unknown type")
	}

	re := regexp.MustCompile(`^(?:hxxp|nsfw|evil|meow)`)

	output = re.ReplaceAllString(output, "http")

	output = strings.Replace(output, "[.]", ".", -1)
	output = strings.Replace(output, "(.)", ".", -1)
	output = strings.Replace(output, "(dot)", ".", -1)

	return output, nil

}
