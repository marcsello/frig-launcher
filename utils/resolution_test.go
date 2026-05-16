package utils

import "testing"

func TestParseResolution(t *testing.T) {
	type args struct {
		res string
	}
	tests := []struct {
		name    string
		args    args
		wantW   int
		wantH   int
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				res: "1024x768",
			},
			wantW:   1024,
			wantH:   768,
			wantErr: false,
		},
		{
			name: "happy2",
			args: args{
				res: "1x1",
			},
			wantW:   1,
			wantH:   1,
			wantErr: false,
		},
		{
			name: "error__negative",
			args: args{
				res: "1x-1",
			},
			wantW:   0,
			wantH:   0,
			wantErr: true,
		},
		{
			name: "error__negative2",
			args: args{
				res: "-1x-1",
			},
			wantW:   0,
			wantH:   0,
			wantErr: true,
		},
		{
			name: "error__zero",
			args: args{
				res: "0x0",
			},
			wantW:   0,
			wantH:   0,
			wantErr: true,
		},
		{
			name: "error__too_many_x",
			args: args{
				res: "1024x768x32",
			},
			wantW:   0,
			wantH:   0,
			wantErr: true,
		},
		{
			name: "error__garbage",
			args: args{
				res: "asdNJHUFADS",
			},
			wantW:   0,
			wantH:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseResolution(tt.args.res)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResolution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantW {
				t.Errorf("ParseResolution() got = %v, want %v", got, tt.wantW)
			}
			if got1 != tt.wantH {
				t.Errorf("ParseResolution() got1 = %v, want %v", got1, tt.wantH)
			}
		})
	}
}
