package defang

import (
	"fmt"
	"net"
	"testing"
)

func ExampleDefangIPv4() {
	defangedIP, _ := IPv4("8.8.8.8")
	fmt.Println(defangedIP)
	// Output: 8.8.8[.]8
}

func TestDefangIPv4(t *testing.T) {
	type args struct {
		ip interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"test string input", args{"4.4.4.4"}, "4.4.4[.]4", false},
		{"test net.IP input", args{net.ParseIP("4.4.4.4")}, "4.4.4[.]4", false},
		{"test unsupported type", args{false}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPv4(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv4() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}
