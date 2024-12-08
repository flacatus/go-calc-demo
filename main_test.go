package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/liornabat-sealights/go-calc-demo/lib/types"
	main2 "github.com/liornabat-sealights/go-calc-demo/pkg/calc"
	"github.com/liornabat-sealights/go-calc-demo/pkg/server"
)

var restyClient = resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

func NewRestRequest() *resty.Request {
	return restyClient.NewRequest()
}

// TestMain is a sample to run an endpoint test
func TestMain(m *testing.M) {
	shuttingDown := make(chan struct{})
	go main()
	time.Sleep(2 * time.Second)
	code := m.Run()
	close(shuttingDown)
	os.Exit(code)
}

func call(path string, a, b string) (*types.ResultResponse, error) {
	time.Sleep(4 * time.Second)
	result := &types.ResultResponse{}
	resp, err := NewRestRequest().SetQueryParams(map[string]string{
		"a": a,
		"b": b,
	}).SetResult(result).
		Get("http://localhost:8080" + path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response: %s", resp.String())
	}
	return result, nil
}
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
			name: "TestAdd",
			args: args{
				a: 1,
				b: 2,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main2.Add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
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
			name: "TestDivide",
			args: args{
				a: 1,
				b: 2,
			},
			want: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main2.Divide(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Divide() = %v, want %v", got, tt.want)
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
			name: "TestMultiply",
			args: args{
				a: 1,
				b: 2,
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main2.Multiply(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResultResponse(t *testing.T) {
	tests := []struct {
		name string
		want *types.ResultResponse
	}{
		{
			name: "TestNewResultResponse",
			want: &types.ResultResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := types.NewResultResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResultResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		want *server.Server
	}{
		{
			name: "TestNewServer",
			want: &server.Server{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := server.NewServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultResponse_SetResult(t *testing.T) {
	type fields struct {
		ValueA float64
		ValueB float64
		Result float64
	}
	type args struct {
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *types.ResultResponse
	}{
		{
			name: "TestResultResponse_SetResult",
			fields: fields{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			args: args{
				value: 4,
			},
			want: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &types.ResultResponse{
				ValueA: tt.fields.ValueA,
				ValueB: tt.fields.ValueB,
				Result: tt.fields.Result,
			}
			if got := r.SetResult(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultResponse_SetValues(t *testing.T) {
	type fields struct {
		ValueA float64
		ValueB float64
		Result float64
	}
	type args struct {
		valueA float64
		valueB float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *types.ResultResponse
	}{
		{
			name: "TestResultResponse_SetValues",
			fields: fields{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			args: args{
				valueA: 4,
			},
			want: &types.ResultResponse{
				ValueA: 4,
				ValueB: 0,
				Result: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &types.ResultResponse{
				ValueA: tt.fields.ValueA,
				ValueB: tt.fields.ValueB,
				Result: tt.fields.Result,
			}
			if got := r.SetValues(tt.args.valueA, tt.args.valueB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_handel_add(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		wantResult *types.ResultResponse
		wantErr    bool
	}{
		{
			name: "TestServer_handelAdd",
			a:    "1",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: false,
		},
		{
			name: "TestServer_handelAdd- bad request",
			a:    "bad",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
		{
			name: "TestServer_handelAdd- bad request 2",
			a:    "1",
			b:    "bad",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
		{
			name: "TestServer_handelAdd- bad request 3",
			a:    "1",
			b:    "",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := call("/add", tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.Result != tt.wantResult.Result {
				t.Errorf("handelAdd() error = %v, wantErr %v", resp.Result, tt.wantResult.Result)
				return
			}
		})
	}
}
func TestServer_handel_sub(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		wantResult *types.ResultResponse
		wantErr    bool
	}{
		{
			name: "TestServer_handelsub",
			a:    "1",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: -1,
			},
			wantErr: false,
		},
		{
			name: "TestServer_handelsub- bad request",
			a:    "bad",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := call("/sub", tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.Result != tt.wantResult.Result {
				t.Errorf("handelAdd() error = %v, wantErr %v", resp.Result, tt.wantResult.Result)
				return
			}
		})
	}
}
func TestServer_handel_mul(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		wantResult *types.ResultResponse
		wantErr    bool
	}{
		{
			name: "TestServer_handelmul",
			a:    "1",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 2,
			},
			wantErr: false,
		},
		{
			name: "TestServer_handelmul- bad request",
			a:    "bad",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := call("/mul", tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.Result != tt.wantResult.Result {
				t.Errorf("handelAdd() error = %v, wantErr %v", resp.Result, tt.wantResult.Result)
				return
			}
		})
	}
}
func TestServer_handel_div(t *testing.T) {
	tests := []struct {
		name       string
		a          string
		b          string
		wantResult *types.ResultResponse
		wantErr    bool
	}{
		{
			name: "TestServer_handeldiv",
			a:    "1",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 0.5,
			},
			wantErr: false,
		},
		{
			name: "TestServer_handeldiv- bad request",
			a:    "bad",
			b:    "2",
			wantResult: &types.ResultResponse{
				ValueA: 1,
				ValueB: 2,
				Result: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := call("/div", tt.a, tt.b)
			if tt.wantErr {
				if err == nil {
					t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("handelAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp.Result != tt.wantResult.Result {
				t.Errorf("handelAdd() error = %v, wantErr %v", resp.Result, tt.wantResult.Result)
				return
			}
		})
	}
}

//
//	func TestServer_handelDiv(t *testing.T) {
//		type fields struct {
//			echoWebServer *echo.Echo
//		}
//		type args struct {
//			c echo.Context
//		}
//		tests := []struct {
//			name    string
//			fields  fields
//			args    args
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				s := &Server{
//					echoWebServer: tt.fields.echoWebServer,
//				}
//				if err := s.handelDiv(tt.args.c); (err != nil) != tt.wantErr {
//					t.Errorf("handelDiv() error = %v, wantErr %v", err, tt.wantErr)
//				}
//			})
//		}
//	}
//
//	func TestServer_handelMul(t *testing.T) {
//		type fields struct {
//			echoWebServer *echo.Echo
//		}
//		type args struct {
//			c echo.Context
//		}
//		tests := []struct {
//			name    string
//			fields  fields
//			args    args
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				s := &Server{
//					echoWebServer: tt.fields.echoWebServer,
//				}
//				if err := s.handelMul(tt.args.c); (err != nil) != tt.wantErr {
//					t.Errorf("handelMul() error = %v, wantErr %v", err, tt.wantErr)
//				}
//			})
//		}
//	}
//
//	func TestServer_handelSub(t *testing.T) {
//		type fields struct {
//			echoWebServer *echo.Echo
//		}
//		type args struct {
//			c echo.Context
//		}
//		tests := []struct {
//			name    string
//			fields  fields
//			args    args
//			wantErr bool
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				s := &Server{
//					echoWebServer: tt.fields.echoWebServer,
//				}
//				if err := s.handelSub(tt.args.c); (err != nil) != tt.wantErr {
//					t.Errorf("handelSub() error = %v, wantErr %v", err, tt.wantErr)
//				}
//			})
//		}
//	}

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
			name: "TestSubtract",
			args: args{
				a: 10,
				b: 5,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := main2.Subtract(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Subtract() = %v, want %v", got, tt.want)
			}
		})
	}
}
