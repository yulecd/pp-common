package trace

import (
	"net"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	got := GetLocalIP()
	if len(got) == 0 {
		t.Errorf("GetLocalIP fail")
	}
}

func TestGetLocalIPDotFormat(t *testing.T) {
	got := GetLocalIPDotFormat()
	if len(got) == 0 {
		t.Errorf("GetLocalIP fail")
	}
}

func Test_isIPUseful(t *testing.T) {
	tests := []struct {
		args net.IP
		want bool
	}{
		{
			[]byte("10.0.0.1"),
			false,
		},
		{
			[]byte("172.16.0.10"),
			false,
		},
		{
			[]byte("192.168.0.20"),
			false,
		},
		{
			[]byte("169.254.0.30"),
			false,
		},
		{
			[]byte("224.0.0.40"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(string(tt.args), func(t *testing.T) {
			if got := isIPUseful(tt.args); got != tt.want {
				t.Errorf("isIPUseful('%v') = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
