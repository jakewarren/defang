package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jakewarren/defang"
	"github.com/mingrammer/commonregex"
	"github.com/spf13/pflag"
)

var version string

type config struct {
	refang bool
}

type app struct {
	Config config
	input  io.Reader
}

func main() {

	d := app{}

	pflag.BoolVarP(&d.Config.refang, "refang", "r", false, "refang IOCs")
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
	clipboard.WriteAll(output)

}

func (d app) defangIOCs() (output string) {
	scanner := bufio.NewScanner(d.input)
	for scanner.Scan() {
		text := scanner.Text()

		// process email addresses
		emails := commonregex.Emails(text)

		for _, e := range emails {

			address := strings.Split(e, "@")

			u, _ := defang.URL(address[1])

			defangedEmail := address[0] + "@" + u

			//remove the email to prevent it from being processed again as a URL
			text = strings.Replace(text, e, "", -1)

			output += defangedEmail + "\n"
		}

		// process links
		links := commonregex.Links(text)

		for _, l := range links {
			u, _ := defang.URL(l)

			output += u + "\n"
		}

	}
	return
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
	if fi.Mode()&os.ModeNamedPipe != 0 { // check if STDIN is attached
		f = os.Stdin
	} else if len(pflag.Args()) > 0 {
		// check the first argument to see if it's a file
		if fileExists(pflag.Arg(0)) {
			f, err = os.Open(pflag.Arg(0))
			if err != nil {
				panic(err)
			}
		} else { // if the user did not pass a file then process the arguments
			f = strings.NewReader(strings.Join(pflag.Args(), "\n"))
		}

	} else {
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
