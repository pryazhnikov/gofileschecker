package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParameters(t *testing.T) {
	// Save original args and restore them after the test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	tests := []struct {
		name    string
		args    []string
		want    *runParameters
		wantErr bool
	}{
		{
			name:    "no arguments at all",
			args:    []string{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "program name only",
			args:    []string{"prog"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid parameters",
			args: []string{"prog", "-path", "/test/path", "-debug", "-fullpath", "-skipempty"},
			want: &runParameters{
				paths:          []string{"/test/path"},
				debug:          true,
				fullFilePath:   true,
				skipEmptyFiles: true,
			},
			wantErr: false,
		},
		{
			name:    "missing path parameter",
			args:    []string{"prog"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only path parameter",
			args: []string{"prog", "-path", "/test/path"},
			want: &runParameters{
				paths:          []string{"/test/path"},
				debug:          false,
				fullFilePath:   false,
				skipEmptyFiles: false,
			},
			wantErr: false,
		},
		{
			name: "multiple path parameters",
			args: []string{"prog", "-path", "/test/path1", "-path", "/test/path2"},
			want: &runParameters{
				paths:          []string{"/test/path1", "/test/path2"},
				debug:          false,
				fullFilePath:   false,
				skipEmptyFiles: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := newRunParametersParser()
			got, err := parser.Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err, "An error is expected")
				return
			}

			assert.NoError(t, err, "parseParameters() should not return an error")
			assert.Equal(t, tt.want.paths, got.paths, "wrong paths value")
			assert.Equal(t, tt.want.debug, got.debug, "wrong value of debug flag")
			assert.Equal(t, tt.want.fullFilePath, got.fullFilePath, "wrong value of fullFilePath flag")
			assert.Equal(t, tt.want.skipEmptyFiles, got.skipEmptyFiles, "wrong value of skipEmptyFiles flag")
		})
	}
}
