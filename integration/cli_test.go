// nolint:scopelint,gosec
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"reflect"
)

var update = flag.Bool("update", false, "update golden files")

var binaryName = "defang"

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), fixture)
}

func writeFixture(t *testing.T, fixture string, content []byte) {
	err := ioutil.WriteFile(fixturePath(t, fixture), content, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestCliArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		fixture string
	}{
		{"no arguments", []string{}, "no-args.golden"},
		{"one argument", []string{"google.com"}, "one-argument.golden"},
		{"file", []string{"testdata/blob.txt"}, "file.golden"},
		{"suricata dns logs", []string{"testdata/suricata_dns_logs.txt"}, "suricata_dns.golden"},
		{"refang file", []string{"-r", "integration/file.golden"}, "refang-file.golden"},
		{"single IP argument", []string{"78.56.216.169"}, "single-ip-argument.golden"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(path.Join(dir, "bin", binaryName), tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("debug: dir: %s\n", dir)
				fmt.Printf("debug: cmd: %s\n", path.Join(dir, "bin", binaryName))
				fmt.Printf("debug: args: %v\n", tt.args)
				fmt.Printf("debug: output: %s\n", output)
				fmt.Printf("debug: error: %s\n", err)
				t.Fatal(err)
			}

			if *update {
				writeFixture(t, tt.fixture, output)
			}

			actual := string(output)

			expected := loadFixture(t, tt.fixture)

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("actual = %s, expected = %s", actual, expected)
			}
		})
	}
}

func TestMain(m *testing.M) {
	err := os.Chdir("..")
	if err != nil {
		fmt.Printf("could not change dir: %v", err)
		os.Exit(1)
	}
	make := exec.Command("make", "build")
	err = make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
