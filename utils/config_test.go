package utils

import (
	"reflect"
	"testing"
)

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
		{"simple", args{path: "../test"}, &Conf{"file:../test/test.db?cache=shared", "8080", "key"}, false},
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
