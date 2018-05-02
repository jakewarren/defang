# defang

[![](https://godoc.org/github.com/jakewarren/defang?status.svg)](http://godoc.org/github.com/jakewarren/defang) 
[![Build Status](https://travis-ci.org/jakewarren/defang.svg?branch=master)](https://travis-ci.org/jakewarren/defang/)
[![GitHub release](http://img.shields.io/github/release/jakewarren/defang.svg?style=flat-square)](https://github.com/jakewarren/defang/releases])
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/defang/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/jakewarren/defang)](https://goreportcard.com/report/github.com/jakewarren/defang)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)

> Defangs and refangs malicious IOCs

## Install
### Option 1: Binary

Download the latest release from [https://github.com/jakewarren/defang/releases/latest](https://github.com/jakewarren/defang/releases/latest)

### Option 2: From source

```
go get github.com/jakewarren/defang/...
```

## Usage
### As a binary:

`defang` accepts input from STDIN, a string argument, or a file argument

```
❯ echo -e "malware.evil.com\napp.malicious.com" | defang
malware.evil[.]com
app.malicious[.]com

❯ defang "malware.evil.com"
malware.evil[.]com

❯ defang /tmp/urls.txt
app.malicious[.]com

❯ defang --refang "hXXp://evil.example[.]com/malicious.php"
http://evil.example.com/malicious.php

```

### As a library:

```
import (
	"fmt"

	"github.com/jakewarren/defang"
)

func main() {
	u, _ := defang.URL("https://malware.evil.com")
	fmt.Println(u)
	// Output: hxxps://malware.evil[.]com

	u, _ = defang.URLWithMask("https://malware.evil.com", defang.Meow)
	fmt.Println(u)
	// Output: meows://malware.evil[.]com

	defangedIP, _ := defang.IPv4("8.8.8.8")
	fmt.Println(defangedIP)
	// Output: 8.8.8[.]8
}
```

## Changes

All notable changes to this project will be documented in the [changelog].

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

## License

MIT © 2018 Jake Warren

[changelog]: https://github.com/jakewarren/defang/blob/master/CHANGELOG.md
