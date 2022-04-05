package mdexec

import (
	"runtime"
	"testing"
)

func Test_getDefaultExecutor(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
		wantErr bool
		skip    bool
	}{
		{
			"simple command works as expected",
			"echo 'test'",
			"test\n",
			false,
			runtime.GOOS == "windows",
		},
		{
			"simple command wrapped with bash",
			`bash -c "echo 'test'"`,
			"test\n",
			false,
			runtime.GOOS == "windows",
		},
		{
			"pipes work",
			`bash -c "echo 'test' | rev"`,
			"tset\n",
			false,
			runtime.GOOS == "windows",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.name)
			}
			executor, err := getDefaultExecutor()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDefaultExecutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _, err := executor(tt.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDefaultExecutor()(%s) error = %v, wantErr %v", tt.command, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}
