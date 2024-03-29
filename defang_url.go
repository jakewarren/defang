package defang

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jakewarren/tldomains"
	"github.com/pkg/errors"
)

// Mask allows the user to apply a different defanging scheme if they are in a fun mood
type Mask int

const (
	// Hxxp applies the default URL scheme of hxxp
	Hxxp Mask = iota
	// Nsfw marks the URL scheme as meow
	Nsfw
	// Meow marks the URL scheme as meow
	Meow
	// Evil marks the URL scheme as evile
	Evil
)

// URL defangs an IPv4 address
func URL(rawURL interface{}) (string, error) {
	var (
		input              string
		output             string
		inputWithoutScheme bool // indicates the url was provided without an URL scheme
	)

	switch u := rawURL.(type) {
	case string:
		input = u
	case url.URL:
		input = u.String()
	default:
		return "", errors.New("unknown type")
	}

	// if the url doesn't have a scheme, add a temporary one so net/url can parse it properly
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
		inputWithoutScheme = true
	}

	u, err := url.Parse(input)
	if err != nil {
		return "", errors.Wrap(err, "error parsing URL (url.Parse)")
	}

	cache := os.TempDir() + "/tld.cache"
	extract, err := tldomains.New(cache)
	if err != nil {
		return "", errors.Wrap(err, "error initializing tldextract")
	}

	var defangedScheme string
	if u.Scheme != "" {
		defangedScheme = strings.Replace(u.Scheme, "http", "hxxp", -1)
	}

	host := extract.Parse(u.Host)
	var defangedHost string
	if len(host.Subdomain) > 0 {
		defangedHost = host.Subdomain + "."
	}

	// if the tld has a period, defang there
	switch {
	case govalidator.IsIPv4(u.Host):
		ip, ipErr := IPv4(u.Host)
		if ipErr != nil {
			return "", errors.New("error defanging IPv4 URL")
		}

		defangedHost = ip
	case strings.Contains(host.Suffix, "."):
		defangedHost += host.Root + "[.]" + strings.Replace(host.Suffix, ".", "[.]", -1)
	case len(host.Suffix) > 0:
		defangedHost += host.Root + "[.]" + host.Suffix
	default:
		return "", errors.New("error defanging URL")

	}

	if len(defangedScheme) > 0 && !inputWithoutScheme {
		output = defangedScheme + "://" + defangedHost
	} else {
		output = defangedHost
	}

	if len(u.Path) > 0 && len(u.Scheme) > 0 {
		output += u.EscapedPath()
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

	if m == 0 {
		return output, nil
	}

	re := regexp.MustCompile(`^hxxp`)

	return re.ReplaceAllString(output, maskStr), nil
}
