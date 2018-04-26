package defang

import (
	"net"
	"regexp"

	"github.com/pkg/errors"
)

// IPv4 defangs an IPv4 address
func IPv4(ip interface{}) (string, error) {

	var ipAddress string

	switch ip := ip.(type) {
	case string:
		ipAddress = ip
	case net.IP:
		ipAddress = ip.String()
	default:
		return "", errors.New("unknown type")
	}

	re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3})(\.)(\d{1,3})`)

	return re.ReplaceAllString(ipAddress, `$1[.]$3`), nil

}
