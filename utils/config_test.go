package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// before tests

	code := m.Run()

	// after tests

	os.Exit(code)
}

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Conf
		wantErr    bool
	}{
		{"wrong path", args{path: ".."}, nil, true},
		{"simple", args{path: "../specs"}, &Conf{"file:../specs/test.db?cache=shared", "8080", "key"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotConfig, err := LoadConfig(tt.args.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("LoadConfig() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
