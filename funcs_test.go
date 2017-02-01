package funcmap

import (
	"math/rand"
	"testing"
)

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
			v := ip_math(m, ip)
			if e != v {
				t.Errorf("expect:%v for:%v result:%v == %v", e, m, v, e == v)
			}
		}
	}
}
