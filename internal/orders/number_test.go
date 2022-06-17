package orders

import "testing"

func TestNumber_Valid(t *testing.T) {
	tests := []struct {
		name string
		num  Number
		want bool
	}{
		{
			name: "valid number 1",
			num:  0,
			want: true,
		},
		{
			name: "valid number 2",
			num:  18,
			want: true,
		},
		{
			name: "valid number 3",
			num:  182,
			want: true,
		},
		{
			name: "valid number 4",
			num:  12345678903,
			want: true,
		},
		{
			name: "valid number 5",
			num:  346436439,
			want: true,
		},
		{
			name: "invalid number 1",
			num:  1,
			want: false,
		},
		{
			name: "invalid number 2",
			num:  10,
			want: false,
		},
		{
			name: "invalid number 3",
			num:  181,
			want: false,
		},
		{
			name: "invalid number 4",
			num:  4665303451001,
			want: false,
		},
		{
			name: "invalid number 5",
			num:  4600936181840,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.num.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
