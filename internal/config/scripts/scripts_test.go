package scripts

import (
	"testing"
)

func TestScripts_Validate(t *testing.T) {
	type fields struct {
		UpdateInterval int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errText string
	}{
		{
			name:    "small updateInterval",
			fields:  fields{UpdateInterval: -1},
			wantErr: true,
			errText: "updateInterval must be not less than 0",
		},
		{
			name:    "ok",
			fields:  fields{UpdateInterval: 1},
			wantErr: false,
			errText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Scripts{
				UpdateInterval: tt.fields.UpdateInterval,
			}
			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errText {
				t.Errorf("Validate() error = '%s', wantErrText '%s'", err.Error(), tt.errText)
			}
		})
	}
}
