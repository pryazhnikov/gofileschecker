package parameters

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
		want    *RunParameters
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
			want: &RunParameters{
				Paths:          []string{"/test/path"},
				Debug:          true,
				FullFilePath:   true,
				SkipEmptyFiles: true,
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
			want: &RunParameters{
				Paths:          []string{"/test/path"},
				Debug:          false,
				FullFilePath:   false,
				SkipEmptyFiles: false,
			},
			wantErr: false,
		},
		{
			name: "multiple path parameters",
			args: []string{"prog", "-path", "/test/path1", "-path", "/test/path2"},
			want: &RunParameters{
				Paths:          []string{"/test/path1", "/test/path2"},
				Debug:          false,
				FullFilePath:   false,
				SkipEmptyFiles: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewRunParametersParser()
			got, err := parser.Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err, "An error is expected")
				return
			}

			assert.NoError(t, err, "parseParameters() should not return an error")
			assert.Equal(t, tt.want.Paths, got.Paths, "wrong paths value")
			assert.Equal(t, tt.want.Debug, got.Debug, "wrong value of debug flag")
			assert.Equal(t, tt.want.FullFilePath, got.FullFilePath, "wrong value of fullFilePath flag")
			assert.Equal(t, tt.want.SkipEmptyFiles, got.SkipEmptyFiles, "wrong value of skipEmptyFiles flag")
		})
	}
}
