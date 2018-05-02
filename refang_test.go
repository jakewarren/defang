package defang

import "testing"

func TestRefang(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"simple http", args{"hxxp://google[.]com"}, "http://google.com", false},
		{"simple http 2", args{"hXXp://evil.example[.]com/malicious.php"}, "http://evil.example.com/malicious.php", false},
		{"simple https", args{"hxxps://google[.]com"}, "https://google.com", false},
		{"meow", args{"meows://google[.]com"}, "https://google.com", false},
		{"dots", args{"hxxps://google(.)com"}, "https://google.com", false},
		{"dots word", args{"hxxps://google(dot)com"}, "https://google.com", false},
		{"dots word multiple", args{"hxxps://google(dot)co(dot)uk"}, "https://google.co.uk", false},
		{"dots in tld 1", args{"hxxps://google[.]co[.]uk"}, "https://google.co.uk", false},
		{"dots in tld 2", args{"hxxps://google.co[.]uk"}, "https://google.co.uk", false},
		{"subdomain", args{"hxxps://ftp[.]example[.]com"}, "https://ftp.example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Refang(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refang() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Refang() = %v, want %v", got, tt.want)
			}
		})
	}
}
