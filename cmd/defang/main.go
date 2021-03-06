package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/atotto/clipboard"
	"github.com/jakewarren/defang"
	"github.com/mingrammer/commonregex"
	"github.com/spf13/pflag"
)

var version string
var ipV4WithPortRE = regexp.MustCompile(`(?m)\d+\.\d+\.\d+\.\d+:\d+`)

type config struct {
	refang  bool
	nsfw    bool
	evil    bool
	meow    bool
	extract bool
}

type app struct {
	Config config
	input  io.Reader
}

func main() {
	d := app{}

	pflag.BoolVarP(&d.Config.refang, "refang", "r", false, "refang IOCs")
	pflag.BoolVarP(&d.Config.extract, "extract", "e", false, "extract IOCs without defanging")
	pflag.BoolVarP(&d.Config.nsfw, "nsfw", "n", false, "defang the URL scheme with nsfw")
	pflag.BoolVar(&d.Config.evil, "evil", false, "defang the URL scheme with evil")
	pflag.BoolVarP(&d.Config.meow, "meow", "m", false, "kitty!")
	displayVersion := pflag.BoolP("version", "V", false, "display version")
	displayHelp := pflag.BoolP("help", "h", false, "display help")

	pflag.Parse()

	// override the default usage display
	if *displayHelp {
		displayUsage()
		os.Exit(0)
	}

	if *displayVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	d.input = getInput()

	d.processInput()
}

func (d app) processInput() {
	var output string

	if d.Config.refang {
		output = d.refangIOCs()
	} else {
		output = d.defangIOCs()
	}

	fmt.Print(output)
	_ = clipboard.WriteAll(strings.TrimSuffix(output, "\n"))
}

func (d app) defangIOCs() (output string) {
	scanner := bufio.NewScanner(d.input)

	m := defang.Hxxp

	switch {
	case d.Config.meow:
		m = defang.Meow
	case d.Config.nsfw:
		m = defang.Nsfw
	case d.Config.evil:
		m = defang.Evil
	}

	for scanner.Scan() {
		text := scanner.Text()

		// process email addresses
		emails := commonregex.Emails(text)

		for _, e := range emails {

			// extract emails without defanging
			if d.Config.extract {
				output += e + "\n"
				continue
			}

			address := strings.Split(e, "@")

			u, err := defang.URLWithMask(address[1], m)
			if err != nil {
				continue
			}

			defangedEmail := address[0] + "@" + u

			// remove the email to prevent it from being processed again as a URL
			text = strings.Replace(text, e, "", -1)

			output += defangedEmail + "\n"
		}

		// process links
		links := commonregex.Links(text)

		for _, l := range links {

			if ipV4WithPortRE.MatchString(l) {
				continue
			}

			if govalidator.IsURL(l) && !govalidator.IsIPv4(l) {
				// extract links without defanging
				if d.Config.extract {
					output += l + "\n"
					continue
				}

				u, err := defang.URLWithMask(l, m)
				if err != nil {
					continue
				}

				output += u + "\n"
			}

		}

		// process IPv4 addresses
		ips := commonregex.IPs(text)
		for _, l := range ips {
			if govalidator.IsIPv4(l) {
				// extract IPs without defanging
				if d.Config.extract {
					output += l + "\n"
					continue
				}

				u, err := defang.IPv4(l)
				if err != nil {
					continue
				}

				output += u + "\n"
			}
		}

	}
	return output
}

func (d app) refangIOCs() (output string) {
	scanner := bufio.NewScanner(d.input)
	for scanner.Scan() {
		text := scanner.Text()

		u, _ := defang.Refang(text)
		output += u + "\n"

	}
	return
}

// getInput determines the input source between STDIN, String param or file name
func getInput() io.Reader {
	var f io.Reader
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	switch {
	case fi.Mode()&os.ModeNamedPipe != 0: // check if STDIN is attached
		f = os.Stdin
	case len(pflag.Args()) > 0:
		// check the first argument to see if it's a file
		if fileExists(pflag.Arg(0)) {
			f, err = os.Open(pflag.Arg(0))
			if err != nil {
				panic(err)
			}
		} else { // if the user did not pass a file then process the arguments
			f = strings.NewReader(strings.Join(pflag.Args(), "\n"))
		}
	default:
		displayUsage()
		os.Exit(0)
	}
	return f
}

// fileExists checks if a file or path fileExists
func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err != nil {
		return os.IsExist(err)
	}
	return true
}

// print custom usage instead of the default provided by pflag
func displayUsage() {
	fmt.Printf("Usage: defang [<flags>] [FILE]\n\n")
	fmt.Printf("Optional flags:\n\n")
	pflag.PrintDefaults()
}
