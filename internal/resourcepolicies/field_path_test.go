package resourcepolicies

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestExtractFieldPathAsString(t *testing.T) {
	type args struct {
		obj       unstructured.Unstructured
		fieldPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractFieldPathAsString(tt.args.obj, tt.args.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFieldPathAsString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractFieldPathAsString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
