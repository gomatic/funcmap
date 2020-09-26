package funcmap

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"text/template"
	"time"

	"github.com/gomatic/clock"
	"github.com/stretchr/testify/assert"
)

//
func TestSubstr(t *testing.T) {
	type param struct {
		start, end int
		s, expect  string
	}
	tests := []param{
		{0, 0, "", ""},
		{0, 0, "0123456789abcdef", ""},
		{0, -1, "0123456789abcdef", "0123456789abcde"},
		{0, -2, "0123456789abcdef", "0123456789abcd"},
		{0, -3, "0123456789abcdef", "0123456789abc"},
		{0, -4, "0123456789abcdef", "0123456789ab"},
		{0, -5, "0123456789abcdef", "0123456789a"},
		{0, -6, "0123456789abcdef", "0123456789"},
		{0, -7, "0123456789abcdef", "012345678"},
		{0, -8, "0123456789abcdef", "01234567"},
		{0, -9, "0123456789abcdef", "0123456"},
		{0, -15, "0123456789abcdef", "0"},
		{0, -16, "0123456789abcdef", ""},
		{0, -17, "0123456789abcdef", "0123456789abcde"},
		{1, 0, "0123456789abcdef", "0"},
		{2, -1, "0123456789abcdef", "23456789abcde"},
		{3, -2, "0123456789abcdef", "3456789abcd"},
		{4, -3, "0123456789abcdef", "456789abc"},
		{5, -4, "0123456789abcdef", "56789ab"},
		{6, -5, "0123456789abcdef", "6789a"},
		{7, -6, "0123456789abcdef", "789"},
		{8, -7, "0123456789abcdef", "8"},
		{9, -8, "0123456789abcdef", "8"},
		{10, -9, "0123456789abcdef", "789"},
		{11, -15, "0123456789abcdef", "123456789a"},
		{12, -16, "0123456789abcdef", "0123456789ab"},
		{13, -17, "0123456789abcdef", "de"},
	}
	for _, p := range tests {
		if got := Substr(p.start, p.end, p.s); got != p.expect {
			t.Errorf("for:%+v got:%v", p, got)
		}
	}
}

//
func TestIpMath(t *testing.T) {
	tests := map[string][][]string{
		"0.0.0.0": {
			{"0.0.0.0", "_._._._"},
			{"255.0.0.0", "[-1]._._._"},
			{"0.0.0.0", "[0]._._._"},
			{"1.0.0.0", "[+1]._._._"},
			{"0.255.0.0", "_.[-1]._._"},
			{"0.0.0.0", "_.[0]._._"},
			{"0.1.0.0", "_.[+1]._._"},
			{"0.0.255.0", "_._.[-1]._"},
			{"0.0.0.0", "_._.[0]._"},
			{"0.0.1.0", "_._.[+1]._"},
			{"0.0.0.255", "_._._.[-1]"},
			{"0.0.0.0", "_._._.[0]"},
			{"0.0.0.1", "_._._.[+1]"},
		},
		"255.255.255.255": {
			{"255.255.255.255", "_._._._"},
			{"254.255.255.255", "[-1]._._._"},
			{"0.255.255.255", "[0]._._._"},
			{"0.255.255.255", "[+1]._._._"},
			{"255.254.255.255", "_.[-1]._._"},
			{"255.0.255.255", "_.[0]._._"},
			{"255.0.255.255", "_.[+1]._._"},
			{"255.255.254.255", "_._.[-1]._"},
			{"255.255.0.255", "_._.[0]._"},
			{"255.255.0.255", "_._.[+1]._"},
			{"255.255.255.254", "_._._.[-1]"},
			{"255.255.255.0", "_._._.[0]"},
			{"255.255.255.0", "_._._.[+1]"},
			{"7.255.255.255", "[-1,/2,%10]._._._"},
			{"1.255.255.255", "[R]._._._"},
			{"5.255.255.255", "[+2,*5,%10]._._._"},
			{"255.7.255.255", "_.[-1,/2,%10]._._"},
			{"255.192.255.255", "_.[R]._._"},
			{"255.5.255.255", "_.[+2,*5,%10]._._"},
			{"255.255.7.255", "_._.[-1,/2,%10]._"},
			{"255.255.115.255", "_._.[R]._"},
			{"255.255.5.255", "_._.[+2,*5,%10]._"},
			{"255.255.255.7", "_._._.[-1,/2,%10]"},
			{"255.255.255.98", "_._._.[R]"},
			{"255.255.255.5", "_._._.[+2,*5,%10]"},
			{"3.255.255.255", "[+R,*R,%R]._._._"},
			{"255.11.255.255", "_.[+R,*R,%R]._._"},
			{"255.255.38.255", "_._.[+R,*R,%R]._"},
			{"255.255.255.40", "_._._.[+R,*R,%R]"},
		},
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff": {
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "_:_:_:_:_:_:_:_"},
			{"fffe:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[-1]:_:_:_:_:_:_:_"},
			{"0000:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[0]:_:_:_:_:_:_:_"},
			{"0000:ffff:ffff:ffff:ffff:ffff:ffff:ffff", "[+1]:_:_:_:_:_:_:_"},
			{"ffff:fffe:ffff:ffff:ffff:ffff:ffff:ffff", "_:[-1]:_:_:_:_:_:_"},
			{"ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff", "_:[0]:_:_:_:_:_:_"},
			{"ffff:0000:ffff:ffff:ffff:ffff:ffff:ffff", "_:[+1]:_:_:_:_:_:_"},
			{"ffff:ffff:fffe:ffff:ffff:ffff:ffff:ffff", "_:_:[-1]:_:_:_:_:_"},
			{"ffff:ffff:0000:ffff:ffff:ffff:ffff:ffff", "_:_:[0]:_:_:_:_:_"},
			{"ffff:ffff:0000:ffff:ffff:ffff:ffff:ffff", "_:_:[+1]:_:_:_:_:_"},
			{"ffff:ffff:ffff:fffe:ffff:ffff:ffff:ffff", "_:_:_:[-1]:_:_:_:_"},
			{"ffff:ffff:ffff:0000:ffff:ffff:ffff:ffff", "_:_:_:[0]:_:_:_:_"},
			{"ffff:ffff:ffff:0000:ffff:ffff:ffff:ffff", "_:_:_:[+1]:_:_:_:_"},
			{"ffff:ffff:ffff:ffff:fffe:ffff:ffff:ffff", "_:_:_:_:[-1]:_:_:_"},
			{"ffff:ffff:ffff:ffff:0000:ffff:ffff:ffff", "_:_:_:_:[0]:_:_:_"},
			{"ffff:ffff:ffff:ffff:0000:ffff:ffff:ffff", "_:_:_:_:[+1]:_:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:fffe:ffff:ffff", "_:_:_:_:_:[-1]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:0000:ffff:ffff", "_:_:_:_:_:[0]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:0000:ffff:ffff", "_:_:_:_:_:[+1]:_:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:fffe:ffff", "_:_:_:_:_:_:[-1]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:0000:ffff", "_:_:_:_:_:_:[0]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:0000:ffff", "_:_:_:_:_:_:[+1]:_"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe", "_:_:_:_:_:_:_:[-1]"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:0000", "_:_:_:_:_:_:_:[0]"},
			{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:0000", "_:_:_:_:_:_:_:[+1]"},
		},
	}

	rand.Seed(0)
	for ip, tests := range tests {
		for _, test := range tests {
			e, m := test[0], test[1]
			v := IPMath(m, ip)
			if e != v {
				t.Errorf("expect:%v for:%v result:%v == %v", e, m, v, e == v)
			}
		}
	}
}

//
func TestAdd(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "0 + 0",
			args: args{
				a: 0,
				b: 0,
			},
			want: 0,
		},
		{
			name: "0 + 1",
			args: args{
				a: 0,
				b: 1,
			},
			want: 1,
		},
		{
			name: "1 + 0",
			args: args{
				a: 0,
				b: 1,
			},
			want: 1,
		},
		{
			name: "0 + -1",
			args: args{
				a: -1,
				b: 0,
			},
			want: -1,
		},
		{
			name: "-1 - 0",
			args: args{
				a: 0,
				b: -1,
			},
			want: -1,
		},
		{
			name: "1 + 1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 2,
		},
		{
			name: "-1 + -1",
			args: args{
				a: -1,
				b: -1,
			},
			want: -2,
		},
		{
			name: "MaxInt64 + MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: math.MaxInt64,
			},
			want: -2,
		},
		{
			name: "MaxInt64 + 0",
			args: args{
				a: math.MaxInt64,
				b: 0,
			},
			want: math.MaxInt64,
		},
		{
			name: "0 + MaxInt64",
			args: args{
				a: 0,
				b: math.MaxInt64,
			},
			want: math.MaxInt64,
		},
		{
			name: "0 + -MaxInt64",
			args: args{
				a: 0,
				b: -math.MaxInt64,
			},
			want: -math.MaxInt64,
		},
		{
			name: "-MaxInt64 + -MaxInt64",
			args: args{
				a: -math.MaxInt64,
				b: -math.MaxInt64,
			},
			want: 2,
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

//
func TestBasename(t *testing.T) {
	type args struct {
		path       string
		extensions []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "path/.remove",
			args: args{
				path:       "/a/b/c/base.remove",
				extensions: []string{"remove"},
			},
			want: "base",
		},
		{
			name: "path/.keep",
			args: args{
				path:       "/a/b/c/base.keep",
				extensions: []string{"remove"},
			},
			want: "base.keep",
		},
		{
			name: "file.keep",
			args: args{
				path:       "base.keep",
				extensions: []string{"remove"},
			},
			want: "base.keep",
		},
		{
			name: "file.remove",
			args: args{
				path:       "file.keep",
				extensions: []string{"remove"},
			},
			want: "file.keep",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Basename(tt.args.path, tt.args.extensions...); got != tt.want {
				t.Errorf("Basename() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestCIDRNext(t *testing.T) {
	type args struct {
		cidr   uint8
		lowest int8
		count  int8
		inc    int8
		addr   []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CIDRNext(tt.args.cidr, tt.args.lowest, tt.args.count, tt.args.inc, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CIDRNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestCleanse(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cleanse(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cleanse() = %v, want %v", got(""), tt.want)
			}
		})
	}
}

//
func TestCleanser(t *testing.T) {
	type args struct {
		r string
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cleanser(tt.args.r, tt.args.s); got != tt.want {
				t.Errorf("Cleanser() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestDebug(t *testing.T) {
	type args struct {
		any []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Debug(tt.args.any...); got != tt.want {
				t.Errorf("Debug() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestDebugger(t *testing.T) {
	tests := []struct {
		name   string
		wantD  bool
		wantDt bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, dt := Debugger()
			if !reflect.DeepEqual(d, tt.wantD) {
				t.Errorf("Debugger() got = %v, want %v", d(), tt.wantD)
			}
			if !reflect.DeepEqual(dt, tt.wantDt) {
				t.Errorf("Debugger() got1 = %v, want %v", dt(), tt.wantDt)
			}
		})
	}
}

//
func TestDecToInt(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecToInt(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestDiv(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "0 / 0",
			args: args{
				a: 0,
				b: 0,
			},
			want: 0,
		},
		{
			name: "1 / 1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 1,
		},
		{
			name: "-1 / -1",
			args: args{
				a: -1,
				b: -1,
			},
			want: 1,
		},
		{
			name: "MaxInt64 / MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: math.MaxInt64,
			},
			want: 1,
		},
		{
			name: "MaxInt64 / 0",
			args: args{
				a: math.MaxInt64,
				b: 0,
			},
			want: 0,
		},
		{
			name: "0 / MaxInt64",
			args: args{
				a: 0,
				b: math.MaxInt64,
			},
			want: 0,
		},
		{
			name: "0 / -MaxInt64",
			args: args{
				a: 0,
				b: -math.MaxInt64,
			},
			want: 0,
		},
		{
			name: "-MaxInt64 / -MaxInt64",
			args: args{
				a: -math.MaxInt64,
				b: -math.MaxInt64,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SafeDiv(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestEnvironment(t *testing.T) {
	type args struct {
		n string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Environment(tt.args.n); got != tt.want {
				t.Errorf("Environment() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestFromInt(t *testing.T) {
	type args struct {
		format string
		arr    []int64
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt(tt.args.format, tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestHexToInt(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexToInt(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP4Add(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint8
		count  uint8
		inc    int8
		addr   []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP4Add(tt.args.group, tt.args.lowest, tt.args.count, tt.args.inc, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IP4Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP4Inc(t *testing.T) {
	type args struct {
		group uint8
		inc   int8
		addr  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP4Inc(tt.args.group, tt.args.inc, tt.args.addr); got != tt.want {
				t.Errorf("IP4Inc() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP4Join(t *testing.T) {
	type args struct {
		addr []int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP4Join(tt.args.addr); got != tt.want {
				t.Errorf("IP4Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP4Next(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint8
		count  uint8
		addr   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP4Next(tt.args.group, tt.args.lowest, tt.args.count, tt.args.addr); got != tt.want {
				t.Errorf("IP4Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP4Prev(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint8
		count  uint8
		addr   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP4Prev(tt.args.group, tt.args.lowest, tt.args.count, tt.args.addr); got != tt.want {
				t.Errorf("IP4Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP6Add(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint16
		count  uint16
		inc    int16
		addr   []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP6Add(tt.args.group, tt.args.lowest, tt.args.count, tt.args.inc, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IP6Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP6Inc(t *testing.T) {
	type args struct {
		group uint8
		inc   int16
		addr  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP6Inc(tt.args.group, tt.args.inc, tt.args.addr); got != tt.want {
				t.Errorf("IP6Inc() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP6Join(t *testing.T) {
	type args struct {
		addr []int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP6Join(tt.args.addr); got != tt.want {
				t.Errorf("IP6Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP6Next(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint16
		count  uint16
		addr   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP6Next(tt.args.group, tt.args.lowest, tt.args.count, tt.args.addr); got != tt.want {
				t.Errorf("IP6Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIP6Prev(t *testing.T) {
	type args struct {
		group  uint8
		lowest uint16
		count  uint16
		addr   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IP6Prev(tt.args.group, tt.args.lowest, tt.args.count, tt.args.addr); got != tt.want {
				t.Errorf("IP6Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIPAdd(t *testing.T) {
	type args struct {
		bits   int32
		group  uint8
		lowest uint16
		count  uint16
		inc    int16
		addr   []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPAdd(tt.args.bits, tt.args.group, tt.args.lowest, tt.args.count, tt.args.inc, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIPCalc(t *testing.T) {
	type args struct {
		bits   int32
		lowest int64
		count  int64
		inc    int64
		value  int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPCalc(tt.args.bits, tt.args.lowest, tt.args.count, tt.args.inc, tt.args.value); got != tt.want {
				t.Errorf("IPCalc() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIPInts(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPInts(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPInts() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIPMath(t *testing.T) {
	type args struct {
		math string
		addr string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPMath(tt.args.math, tt.args.addr); got != tt.want {
				t.Errorf("IPMath() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIPSplit(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPSplit(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestIndex(t *testing.T) {
	type args struct {
		i int
		a interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Index(tt.args.i, tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestJoin(t *testing.T) {
	type args struct {
		sep string
		arr []string
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := Join(tt.args.sep, tt.args.arr); gotS != tt.wantS {
				t.Errorf("Join() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

//
func TestKeySequencer(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeySequencer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeySequencer() = %v, want %v", got(tt.name), tt.want)
			}
		})
	}
}

//
func TestMod(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mod(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestMul(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{

		{
			name: "0 * 0",
			args: args{
				a: 0,
				b: 0,
			},
			want: 0,
		},
		{
			name: "1 * 0",
			args: args{
				a: 0,
				b: 1,
			},
			want: 0,
		},
		{
			name: "-1 * 0",
			args: args{
				a: 0,
				b: -1,
			},
			want: 0,
		},
		{
			name: "0 * -1",
			args: args{
				a: -1,
				b: 0,
			},
			want: 0,
		},
		{
			name: "1 * 1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 1,
		},
		{
			name: "-1 * -1",
			args: args{
				a: -1,
				b: -1,
			},
			want: 1,
		},
		{
			name: "MaxInt64 * MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: math.MaxInt64,
			},
			want: 1,
		},
		{
			name: "0 * MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: 0,
			},
			want: 0,
		},
		{
			name: "MaxInt64 * 0",
			args: args{
				a: 0,
				b: math.MaxInt64,
			},
			want: 0,
		},
		{
			name: "-MaxInt64 * 0",
			args: args{
				a: 0,
				b: -math.MaxInt64,
			},
			want: 0,
		},
		{
			name: "-MaxInt64 * -MaxInt64",
			args: args{
				a: -math.MaxInt64,
				b: -math.MaxInt64,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mul(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestNew(t *testing.T) {
	type args struct {
		options []Optional
	}
	tests := []struct {
		name string
		args args
		want template.FuncMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestParseInt(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int64
	}{
		{
			name: "20",
			args: 10,
			want: 20,
		},
		{
			name: "12",
			args: 10,
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if parser := IntParser(tt.args); !reflect.DeepEqual(parser, tt.want) {
				i, err := parser(tt.name)
				assert.NoErrorf(t, err, fmt.Sprintf("parsing: %s", tt.name))
				assert.Equal(t, tt.want, i)
			}
		})
	}
}

//
func TestPause(t *testing.T) {
	type args struct {
		t int64
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pause(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pause() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestRand(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Rand(); got != tt.want {
				t.Errorf("Rand() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestReInitcap(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReInitcap(tt.args); got != tt.want {
				t.Errorf("ReInitcap() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestReReplace(t *testing.T) {
	type args struct {
		n   int
		old string
		new string
		s   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReReplace(tt.args.n, tt.args.old, tt.args.new, tt.args.s); got != tt.want {
				t.Errorf("ReReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestReTrim(t *testing.T) {
	type args struct {
		cut string
		s   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReTrim(tt.args.cut, tt.args.s); got != tt.want {
				t.Errorf("ReTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestReTrimLeft(t *testing.T) {
	type args struct {
		cut string
		s   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReTrimLeft(tt.args.cut, tt.args.s); got != tt.want {
				t.Errorf("ReTrimLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestReTrimRight(t *testing.T) {
	type args struct {
		cut string
		s   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReTrimRight(tt.args.cut, tt.args.s); got != tt.want {
				t.Errorf("ReTrimRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestSequencer(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sequencer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sequencer() = %v, want %v", got(), tt.want)
			}
		})
	}
}

//
func TestSplit(t *testing.T) {
	type args struct {
		sep string
		s   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.args.sep, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestStarter(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Starter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Starter() = %v, want %v", got(), tt.want)
			}
		})
	}
}

//
func TestStep(t *testing.T) {
	type args struct {
		a  int64
		is []int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Step(tt.args.a, tt.args.is...); got != tt.want {
				t.Errorf("Step() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestSub(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "0 - 0",
			args: args{
				a: 0,
				b: 0,
			},
			want: 0,
		},
		{
			name: "1 - 0",
			args: args{
				a: 0,
				b: 1,
			},
			want: 1,
		},
		{
			name: "-1 - 0",
			args: args{
				a: 0,
				b: -1,
			},
			want: -1,
		},
		{
			name: "0 - -1",
			args: args{
				a: -1,
				b: 0,
			},
			want: 1,
		},
		{
			name: "1 - 1",
			args: args{
				a: 1,
				b: 1,
			},
			want: 0,
		},
		{
			name: "-1 - -1",
			args: args{
				a: -1,
				b: -1,
			},
			want: 0,
		},
		{
			name: "MaxInt64 - MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: math.MaxInt64,
			},
			want: 0,
		},
		{
			name: "0 - MaxInt64",
			args: args{
				a: math.MaxInt64,
				b: 0,
			},
			want: -9223372036854775807,
		},
		{
			name: "MaxInt64 - 0",
			args: args{
				a: 0,
				b: math.MaxInt64,
			},
			want: math.MaxInt64,
		},
		{
			name: "-MaxInt64 - 0",
			args: args{
				a: 0,
				b: -math.MaxInt64,
			},
			want: -math.MaxInt64,
		},
		{
			name: "-MaxInt64 - -MaxInt64",
			args: args{
				a: -math.MaxInt64,
				b: -math.MaxInt64,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sub(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestSubstr1(t *testing.T) {
	type args struct {
		start int
		end   int
		s     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substr(tt.args.start, tt.args.end, tt.args.s); got != tt.want {
				t.Errorf("Substr() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestToInt(t *testing.T) {
	type args struct {
		base int
		arr  []string
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.base, tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestWithClock(t *testing.T) {
	type args struct {
		timeFunc clock.TimeFunction
	}
	tests := []struct {
		name string
		args args
		want Optional
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithClock(tt.args.timeFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestWithMaps(t *testing.T) {
	tests := []struct {
		name string
		args []template.FuncMap
		opts Optional
		want template.FuncMap
	}{
		{
			name: "empty",
			args: []template.FuncMap{},
			want: template.FuncMap{},
		},
		{
			name: "simple",
			args: []template.FuncMap{{"add": Add, "sub": Sub}},
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
		{
			name: "duplicates",
			args: []template.FuncMap{{"add": Add, "sub": Sub}, {"add": Add, "sub": Sub}},
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
		{
			name: "replacements",
			args: []template.FuncMap{{"add": Add, "sub": Sub}, {"add": Div, "sub": Mul}},
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
		{
			name: "ignore nil",
			args: []template.FuncMap{{"add": nil, "sub": nil}, {"add": Add, "sub": Sub}},
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
		{
			name: "duplicates overrides",
			args: []template.FuncMap{{"add": Add, "sub": Sub}, {"add": Add, "sub": Sub}},
			opts: WithRightmostOverrides(),
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
		{
			name: "replacements overrides",
			args: []template.FuncMap{{"add": Add, "sub": Sub}, {"add": Div, "sub": Mul}},
			opts: WithRightmostOverrides(),
			want: template.FuncMap{"add": Div, "sub": Mul},
		},
		{
			name: "ignore nil overrides",
			args: []template.FuncMap{{"add": nil, "sub": nil}, {"add": Add, "sub": Sub}},
			opts: WithRightmostOverrides(),
			want: template.FuncMap{"add": Add, "sub": Sub},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(WithMaps(tt.args...), tt.opts)
			for k, _ := range tt.want {
				got := reflect.ValueOf(got[k])
				expect := reflect.ValueOf(tt.want[k])
				assert.Equal(t, &got, &expect)
			}
		})
	}
}
