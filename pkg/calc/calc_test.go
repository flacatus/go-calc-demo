package calc

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Add positive numbers",
			args: args{a: 4.5, b: 5.5},
			want: 10,
		},
		{
			name: "Add negative numbers",
			args: args{a: -1.5, b: -2.5},
			want: -4,
		},
		{
			name: "Add positive and negative numbers",
			args: args{a: 5, b: -3},
			want: 2,
		},
		{
			name: "Add zero and number",
			args: args{a: 0, b: 5},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Positive and positive",
			args: args{a: 10.0, b: 5.0},
			want: 5.0,
		},
		{
			name: "Positive and negative",
			args: args{a: 10.0, b: -5.0},
			want: 15.0,
		},
		{
			name: "Negative and positive",
			args: args{a: -10.0, b: 5.0},
			want: -15.0,
		},
		{
			name: "Negative and negative",
			args: args{a: -10.0, b: -5.0},
			want: -5.0,
		},
		{
			name: "Zero subtract positive",
			args: args{a: 0.0, b: 5.0},
			want: -5.0,
		},
		{
			name: "Zero subtract negative",
			args: args{a: 0.0, b: -5.0},
			want: 5.0,
		},
		{
			name: "Zero subtract zero",
			args: args{a: 0.0, b: 0.0},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Subtract(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Multiply positive numbers",
			args: args{a: 4.5, b: 2},
			want: 9,
		},
		{
			name: "Multiply negative numbers",
			args: args{a: -3, b: -2},
			want: 6,
		},
		{
			name: "Multiply positive and negative numbers",
			args: args{a: -3, b: 4},
			want: -12,
		},
		{
			name: "Multiply by zero",
			args: args{a: 0, b: 5},
			want: 0,
		},
		{
			name: "Multiply zero by zero",
			args: args{a: 0, b: 0},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Multiply(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Divide positive numbers",
			args: args{a: 10, b: 2},
			want: 5,
		},
		{
			name: "Divide negative numbers",
			args: args{a: -10, b: -2},
			want: 5,
		},
		{
			name: "Divide positive by negative number",
			args: args{a: 10, b: -2},
			want: -5,
		},
		{
			name: "Divide zero by number",
			args: args{a: 0, b: 2},
			want: 0,
		},
		{
			name: "Divide number by zero",
			args: args{a: 10, b: 0},
			want: math.Inf(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Divide(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}
