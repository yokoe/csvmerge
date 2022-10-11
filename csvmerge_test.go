package csvmerge

import (
	"testing"
)

func TestMergeCSVFiles(t *testing.T) {
	type args struct {
		files []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"No files", args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MergeCSVFiles(tt.args.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeCSVFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
