package available

import (
	"errors"
	"reflect"
	"testing"

	"github.com/openshift/backplane-tools/pkg/tools"
)

type fakeTool struct {
	tools.Tool
	version string
	err     error
}

func (t fakeTool) LatestVersion() (string, error) {
	return t.version, t.err
}

func TestLatestVersions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		toolMap map[string]tools.Tool
		want    map[string]string
		wantErr bool
	}{
		{
			name: "all versions resolve",
			toolMap: map[string]tools.Tool{
				"oc": fakeTool{version: "4.17.0"},
				"yq": fakeTool{version: "v4.44.3"},
			},
			want: map[string]string{
				"oc": "4.17.0",
				"yq": "v4.44.3",
			},
		},
		{
			name: "any lookup failure fails the whole result",
			toolMap: map[string]tools.Tool{
				"oc": fakeTool{version: "4.17.0"},
				"yq": fakeTool{err: errors.New("rate limited")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := LatestVersions(tt.toolMap)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LatestVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LatestVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
