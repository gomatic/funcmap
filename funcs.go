package funcmap

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//
var Map = template.FuncMap{
	"debug":        debug,
	"command_line": commandLine,
	"ip_math":      ip_math,
	"ip4_inc":      ip4_inc,
	"ip4_next":     ip4_next,
	"ip4_prev":     ip4_prev,
	"ip4_add":      ip4_add,
	"ip4_join":     ip4_join,
	"ip6_inc":      ip6_inc,
	"ip6_next":     ip6_next,
	"ip6_prev":     ip6_prev,
	"ip6_add":      ip6_add,
	"ip6_join":     ip6_join,
	"cidr_next":    cidr_next,
	"ip_ints":      ip_ints,
	"ip_split":     ip_split,
	"to_int":       to_int,
	"dec_to_int":   dec_to_int,
	"hex_to_int":   hex_to_int,
	"from_int":     from_int,
	"inc":          step,
	"add":          add,
	"sub":          sub,
	"mul":          mul,
	"div":          div,
	"mod":          mod,
	"rand":         func() int64 { return rand.Int63() },
	"cleanse":      cleanse(`[^[:alpha:]]`),
	"environment":  environment,
	"now":          time.Now,
	"started":      started(),
	"iindex":       index,
	"split":        split,
	"join":         join,
	"substr":       substr,
	"lower":        strings.ToLower,
	"replace":      strings.Replace,
	"trim":         strings.Trim,
	"trim_left":    strings.TrimLeft,
	"trim_right":   strings.TrimRight,
	"upper":        strings.ToUpper,
}

// To report a consistent time through a single template.
var started = func() func() time.Time {
	started := time.Now()
	return func() time.Time { return started }
}

//
func debug(any ...interface{}) string {
	s := make([]string, len(any))
	for i, a := range any {
		s[i] = fmt.Sprintf("%v", a)
	}
	return join(" ", s)
}

//
func step(a int64, is ...int) int64 {
	if len(is) == 0 {
		is = []int{1}
	}
	for _, i := range is {
		a += int64(i)
	}
	return a
}

func add(a, b int64) int64 { return b + a }

// Subtract `a` from `b`
func sub(a, b int64) int64 { return b - a }
func mul(a, b int64) int64 { return b * a }

// `b` modulo `a`
func mod(a, b int64) int64 { return b % a }

// `b` divided by `a`. Returns `0` if `a == 0`.
func div(a, b int64) int64 {
	if a == 0 {
		return 0
	}
	return b / a
}

//
func cleanse(r string) func(string) string {
	re := regexp.MustCompile(r)
	return func(s string) string {
		return re.ReplaceAllString(s, "")
	}
}

func parseInt(base int) func(s string) (int64, error) {
	return func(s string) (int64, error) {
		return strconv.ParseInt(s, base, 64)
	}
}

//
func environment(n string) string {
	v, _ := os.LookupEnv(n)
	return v
}

var (
	parseDec = parseInt(10)
	parseHex = parseInt(16)
)

// TODO increment CIDR
func cidr_next(cidr uint8, lowest, count, inc int8, addr []int64) []int64 {
	return addr
}

func ip_calc(bits int32, lowest, count, inc, value int64) int64 {
	if value < lowest {
		value += int64(bits)
	}
	return (lowest + (value-lowest+inc)%count) % int64(bits)
}

// Given a zero-based, left-to-right IP group index, lowest value, count, and increment,
// increment the group, cyclically.
func ip_add(bits int32, group uint8, lowest, count uint16, inc int16, addr []int64) []int64 {
	if group >= uint8(len(addr)) {
		return addr
	}
	if lowest == 0 && count == 0 {
		addr[group] = (addr[group] + int64(inc)) % int64(bits)
	} else {
		addr[group] = ip_calc(int32(bits), int64(lowest), int64(count), int64(inc), addr[group])
	}
	return addr
}

//
func ip4_inc(group uint8, inc int8, addr string) string {
	return ip4_join(ip4_add(group, 0, 0, inc, ip_ints(addr)))
}

//
func ip4_next(group uint8, lowest, count uint8, addr string) string {
	return ip4_join(ip4_add(group, lowest, count, 1, ip_ints(addr)))
}

//
func ip4_prev(group uint8, lowest, count uint8, addr string) string {
	return ip4_join(ip4_add(group, lowest, count, -1, ip_ints(addr)))
}

// Given a zero-based, left-to-right IP group index, lowest value, count, and increment,
// increment the group, cyclically.
func ip4_add(group uint8, lowest, count uint8, inc int8, addr []int64) []int64 {
	return ip_add(int32(256), group, uint16(lowest), uint16(count), int16(inc), addr)
}

//
func ip6_inc(group uint8, inc int16, addr string) string {
	return ip6_join(ip6_add(group, 0, 0, inc, ip_ints(addr)))
}

//
func ip6_next(group uint8, lowest, count uint16, addr string) string {
	return ip6_join(ip6_add(group, lowest, count, 1, ip_ints(addr)))
}

//
func ip6_prev(group uint8, lowest, count uint16, addr string) string {
	return ip6_join(ip6_add(group, lowest, count, -1, ip_ints(addr)))
}

// given a group, lowest, count, and increment, increment the group, circling around
func ip6_add(group uint8, lowest, count uint16, inc int16, addr []int64) []int64 {
	return ip_add(int32(65536), group, lowest, count, inc, addr)
}

//
func join(sep string, arr []string) (s string) {
	return strings.Join(arr, sep)
}

//
func substr(start, end int, s string) string {
	l := len(s)
	if l == 0 {
		return s
	}
	start, end = start%l, end%l
	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}
	if start > end {
		start, end = end, start
	}
	if start > l || start < 0 || end < 0 {
		return s
	} else if end > l {
		end = l
	}
	return s[start:end]
}

//
func split(sep, s string) []string {
	return strings.Split(s, sep)
}

//
func index(i int, a []int64) int64 {
	if a == nil || i < 0 || i >= len(a) {
		return -1
	}
	return a[i]
}

//
func ip_split(addr string) []string {
	ip_groups := split(".", addr)
	if len(ip_groups) > 1 {
		return ip_groups
	}
	return split(":", addr)
}

//
func ip4_join(addr []int64) string {
	return join(".", from_int("%d", addr))
}

//
func ip6_join(addr []int64) string {
	return join(":", from_int("%04x", addr))
}

//
func ip_ints(addr string) []int64 {
	if ip_groups := split(".", addr); len(ip_groups) > 1 {
		return dec_to_int(ip_groups)
	} else {
		return hex_to_int(strings.Split(":", addr))
	}
}

//
func dec_to_int(arr []string) []int64 {
	return to_int(10, arr)
}

//
func hex_to_int(arr []string) []int64 {
	return to_int(16, arr)
}

//
func to_int(base int, arr []string) []int64 {
	is := make([]int64, len(arr))
	parser := parseInt(base)
	for i, m := range arr {
		p, err := parser(m)
		if err != nil {
			continue
		}
		is[i] = p
	}
	return is
}

//
func from_int(format string, arr []int64) []string {
	ss := make([]string, len(arr))
	for i, m := range arr {
		ss[i] = fmt.Sprintf(format, m)
	}
	return ss
}

// Performs IP math using a simple sequence of operations.
// e.g. _.[+2]._.[+1,%10]
func ip_math(math, addr string) string {
	sep, format, width := ".", "%d", uint(256)
	ip_groups := split(sep, addr)
	th_groups := split(sep, math)
	parser := parseDec
	if len(ip_groups) == 1 {
		parser = parseHex
		sep, format, width = ":", "%04x", uint(65536)
		ip_groups = split(sep, addr)
		th_groups = split(sep, math)
	}
	if len(ip_groups) != len(th_groups) {
		return addr
	}
	ip_values := make([]int64, len(ip_groups))
	for i, m := range ip_groups {
		p, err := parser(m)
		if err != nil {
			continue
		}
		ip_values[i] = p
	}
	for i, m := range th_groups {
		m := m
		lm := len(m)
		if lm < 3 {
			continue
		}
		switch m {
		case "_":
			continue
		}
		if m[0] != '[' || m[lm-1] != ']' {
			continue
		}
		m = m[1: lm-1]
		p := ip_values[i]
		for _, a := range strings.Split(m, ",") {
			a := a
			op := a[0]
			switch op {
			case '+', '-', '*', '/', '%':
				a = a[1:]
			default:
			}

			n := int64(0)
			switch a {
			case "R":
				n = rand.Int63n(int64(width))
			default:
				x, err := parser(a)
				if err != nil {
					continue
				}
				n = x
			}

			switch op {
			case '+':
				p += n
			case '-':
				p -= n
			case '*':
				p *= n
			case '/':
				p /= n
			case '%':
				p %= n
			default:
				p = n
			}
			p %= int64(width)
		}
		ip_groups[i] = fmt.Sprintf(format, uint(p)%width)
	}
	return join(sep, ip_groups)
}

// Reproduce a command line string that reflects a usable command line.
func commandLine() string {

	quoter := func(e string) string {
		if !strings.Contains(e, " ") {
			return e
		}
		p := strings.SplitN(e, "=", 2)
		if strings.Contains(p[0], " ") {
			p[0] = `"` + strings.Replace(p[0], `"`, `\"`, -1) + `"`
		}
		if len(p) == 1 {
			return p[0]
		}
		return p[0] + `="` + strings.Replace(p[1], `"`, `\"`, -1) + `"`
	}
	each := func(s []string) (o []string) {
		o = make([]string, len(s))
		for i, t := range s {
			o[i] = quoter(t)
		}
		return
	}
	return filepath.Base(os.Args[0]) + " " + strings.Join(each(os.Args[1:]), " ")
}
