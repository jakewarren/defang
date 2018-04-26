package defang

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/joeguo/tldextract"
	"github.com/pkg/errors"
)

// Mask allows the user to apply a different defanging scheme if they are in a fun mood
type Mask int

const (
	// Nsfw marks the URL as Nsfw
	Nsfw Mask = iota
	Meow
	Evil
)

// URL defangs an IPv4 address
func URL(rawURL interface{}) (string, error) {

	var (
		input  string
		output string
	)

	switch u := rawURL.(type) {
	case string:
		input = u
	case url.URL:
		input = u.String()
	default:
		return "", errors.New("unknown type")
	}

	u, _ := url.Parse(input)

	cache := os.TempDir() + "/tld.cache"
	extract, err := tldextract.New(cache, false)

	if err != nil {
		return "", errors.Wrap(err, "error parsing URL")
	}

	defangedScheme := strings.Replace(u.Scheme, "http", "hxxp", -1)
	host := extract.Extract(input)
	var defangedHost string
	if len(host.Sub) > 0 {
		defangedHost = host.Sub + "."
	}

	// if the tld has a period, defang there
	if strings.Contains(host.Tld, ".") {
		defangedHost += host.Root + "." + strings.Replace(host.Tld, ".", "[.]", -1)
	} else if len(host.Tld) > 0 {
		defangedHost += host.Root + "[.]" + host.Tld
	} else if govalidator.IsIPv4(host.Root) {
		ip, ipErr := IPv4(host.Root)
		if ipErr != nil {
			return "", errors.New("error defanging IPv4 URL")
		}

		defangedHost += ip
	} else {
		return "", errors.New("error defanging URL")
	}

	if len(defangedScheme) > 0 {
		output = defangedScheme + "://" + defangedHost
	} else {
		output = defangedHost
	}

	if len(u.Path) > 0 && len(u.Scheme) > 0 {
		output += u.Path
	}

	if len(u.RawQuery) > 0 {
		output += "?" + u.RawQuery
	}

	if len(u.Fragment) > 0 {
		output += "#" + u.Fragment
	}

	return output, nil
}

// URLWithMask defangs a URL and applies the specified defanging scheme
func URLWithMask(rawURL interface{}, m Mask) (string, error) {

	output, err := URL(rawURL)
	if err != nil {
		return "", err
	}
	var maskStr string
	switch m {
	case Nsfw:
		maskStr = "nsfw"
	case Meow:
		maskStr = "meow"
	case Evil:
		maskStr = "evil"
	}

	re := regexp.MustCompile(`^hxxp`)

	return re.ReplaceAllString(output, maskStr), nil

}
