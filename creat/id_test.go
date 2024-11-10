package creat

import "testing"

func TestIDCompareTo(t *testing.T) {
	tests := []struct {
		id1, id2 ID
		want     int
	}{
		{
			id1: ID("a"), id2: ID("a"),
			want: 0,
		},
		{
			id1: ID("a"), id2: ID("b"),
			want: -1,
		},
		{
			id1: ID("b"), id2: ID("a"),
			want: 1,
		},
		{
			id1: ID("abc"), id2: ID("abd"),
			want: -1,
		},
		{
			id1: ID("abd"), id2: ID("abc"),
			want: 1,
		},
	}
	for _, test := range tests {
		t.Run(string(test.id1)+"_"+string(test.id2), func(t *testing.T) {
			if got := test.id1.CompareTo(test.id2); got != test.want {
				t.Errorf("CompareTo() = %v, want %v", got, test.want)
			}
		})
	}
}
