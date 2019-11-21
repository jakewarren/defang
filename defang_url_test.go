// nolint:scopelint
package defang

import (
	"fmt"
	"testing"
)

func ExampleURL() {
	d, _ := URL("https://subdomain.bing.co.uk/search?q=testquery#testanchor")
	fmt.Println(d)
	// Output: hxxps://subdomain.bing[.]co[.]uk/search?q=testquery#testanchor
}

func ExampleURLWithMask() {
	d, _ := URLWithMask("https://subdomain.bing.co.uk/search?q=testquery#testanchor", Meow)
	fmt.Println(d)
	// Output: meows://subdomain.bing[.]co[.]uk/search?q=testquery#testanchor
}

func TestURL(t *testing.T) {
	type args struct {
		rawURL interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple http", args{"http://www.google.com"}, "hxxp://www.google[.]com", false},
		{"simple https", args{"https://www.google.com"}, "hxxps://www.google[.]com", false},
		{"multiple subdomains", args{"https://sub.www.google.com"}, "hxxps://sub.www.google[.]com", false},
		{"complex tld", args{"https://www.google.co.uk"}, "hxxps://www.google[.]co[.]uk", false},
		{"retain URL fragment", args{"https://www.google.co.uk/foobar#baz"}, "hxxps://www.google[.]co[.]uk/foobar#baz", false},
		{"IPv4 URL", args{"https://1.2.3.4/foobar"}, "hxxps://1.2.3[.]4/foobar", false},
		{"simple url", args{"google.com"}, "google[.]com", false},
		{"complex URL", args{"example.com/dir/subdir/"}, "example[.]com/dir/subdir/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLWithMask(t *testing.T) {
	type args struct {
		rawURL interface{}
		m      Mask
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{"https://www.google.com", Nsfw}, "nsfws://www.google[.]com", false},
		{"2", args{"https://www.google.com", Evil}, "evils://www.google[.]com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLWithMask(tt.args.rawURL, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLWithMask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLWithMask() = %v, want %v", got, tt.want)
			}
		})
	}
}
