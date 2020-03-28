package function

import (
	"context"
	"testing"
)

func TestIntegrationECR_CheckRepositoryCompliance(t *testing.T) {
	type args struct {
		ctx   context.Context
		image string
	}
	tests := []struct {
		name    string
		c       *Container
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.CheckRepositoryCompliance(tt.args.ctx, tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("Container.CheckRepositoryCompliance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Container.CheckRepositoryCompliance() = %v, want %v", got, tt.want)
			}
		})
	}
}
