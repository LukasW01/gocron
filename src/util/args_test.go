package util

import (
	"os"
	"reflect"
	"testing"
)

func TestSplitArgs(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name      string
		args      []string
		wantFlags []string
		wantExec  []string
	}{
		{
			name:      "No separator",
			args:      []string{"cmd", "--foo", "bar"},
			wantFlags: []string{"cmd", "--foo", "bar"},
			wantExec:  nil,
		},
		{
			name:      "With separator",
			args:      []string{"cmd", "--foo", "bar", "--", "exec", "--opt"},
			wantFlags: []string{"cmd", "--foo", "bar"},
			wantExec:  []string{"exec", "--opt"},
		},
		{
			name:      "Only separator",
			args:      []string{"cmd", "--"},
			wantFlags: []string{"cmd"},
			wantExec:  []string{},
		},
		{
			name:      "Separator at beginning",
			args:      []string{"--", "exec"},
			wantFlags: []string{},
			wantExec:  []string{"exec"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			gotFlags, gotExec := SplitArgs()

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("Flags = %v, want %v", gotFlags, tt.wantFlags)
			}
			if !reflect.DeepEqual(gotExec, tt.wantExec) {
				t.Errorf("ExecArgs = %v, want %v", gotExec, tt.wantExec)
			}
		})
	}
}
