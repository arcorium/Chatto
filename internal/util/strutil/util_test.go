package strutil

import "testing"

func TestIsEmpty(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "non empty",
			args: args{"non-empty string"},
			want: false,
		},
		{
			name: "empty",
			args: args{""},
			want: true,
		},
		{
			name: "non-initialized value",
			args: args{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.str); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
