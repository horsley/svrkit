package svrkit

import (
	"reflect"
	"testing"
)

func TestInSliceChecker(t *testing.T) {
	type args struct {
		slice,
		needle interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 1",
			args: args{
				slice:  []int{5, 7, 9},
				needle: 7,
			},
			want: true,
		},
		{
			name: "test 2",
			args: args{
				slice:  []int{5, 7, 9},
				needle: 8,
			},
			want: false,
		},
		{
			name: "test 3",
			args: args{
				slice:  1,
				needle: 8,
			},
			want: false,
		},
		{
			name: "test 4",
			args: args{
				slice:  []string{"1213", "212", "211dd1"},
				needle: "212",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InSliceChecker(tt.args.slice); got(tt.args.needle) != tt.want {
				t.Errorf("InSliceChecker() = %v, want %v", got(tt.args.needle), tt.want)
			}
		})
	}
}

func TestSliceToList(t *testing.T) {
	type args struct {
		slice    interface{}
		seprator string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				slice: []string{"horsley", "panpan", "kavin"},
			},
			want: "horsley,panpan,kavin",
		},
		{
			name: "test 2",
			args: args{
				slice:    []string{"horsley", "panpan", "kavin"},
				seprator: ";",
			},
			want: "horsley;panpan;kavin",
		},
		{
			name: "test 3",
			args: args{
				slice:    []int{1, 4, 7},
				seprator: ",",
			},
			want: "1,4,7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceToList(tt.args.slice, tt.args.seprator); got != tt.want {
				t.Errorf("SliceToList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSliceUnique(t *testing.T) {
	type args struct {
		src []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test 1",
			args: args{
				src: []int{1, 1, 3, 5, 4, 6, 1, 2, 4, 3},
			},
			want: []int{1, 3, 5, 4, 6, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntSliceUnique(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntSliceUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceUnique(t *testing.T) {
	type args struct {
		src []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test 1",
			args: args{
				src: []string{"11", "1daa", "sda", "1daa", "311fpwe", "1daa", "11"},
			},
			want: []string{"11", "1daa", "sda", "311fpwe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceUnique(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}
