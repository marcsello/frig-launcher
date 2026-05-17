package utils

import "testing"

func TestBoolToInt(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "true",
			args: args{
				b: true,
			},
			want: 1,
		},
		{
			name: "false",
			args: args{
				b: false,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolToInt(tt.args.b); got != tt.want {
				t.Errorf("BoolToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
