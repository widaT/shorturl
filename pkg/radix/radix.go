package radix

import "errors"

var ErrInvalidString = errors.New("radix: invalid string")

type Radix struct {
	charset    []byte
	charsetRef map[byte]int
	offset     func(i int) int
	radix      int
	radixMod   int // quick mod
}

func New(opts ...apply) *Radix {
	r := new(Radix)
	for _, opt := range newOptions(opts...) {
		opt(r)
	}
	return r
}

// ParseInt interprets a string str in the given charset returns the corresponding value i
func (r *Radix) ParseInt(str string) (int, error) {
	val, fac := 0, 1
	for pos := len(str) - 1; pos >= 0; pos-- {
		v, ok := r.charsetRef[str[pos]]
		if !ok {
			return 0, ErrInvalidString
		}
		if r.radixMod == 0 {
			u := v - r.offset(pos)%r.radix
			if u < 0 {
				u += r.radix
			}
			val += fac * (u % r.radix)
		} else {
			val += fac * ((v - r.offset(pos)) & r.radixMod)
		}
		fac *= r.radix
	}
	return val, nil
}

// FormatInt returns the string representation of i in the given charset
func (r *Radix) FormatInt(num int, b []byte) {
	for pos := len(b) - 1; pos >= 0; pos-- {
		b[pos] = r.Format(num, pos)
		num /= r.radix
	}
}

// Format returns the byte representation of i in the given charset
func (r *Radix) Format(num int, pos int) byte {
	if r.radixMod == 0 {
		return r.charset[(num+r.offset(pos))%r.radix]
	} else {
		return r.charset[(num+r.offset(pos))&r.radixMod]
	}
}

// Itoa is equivalent to FormatInt(num, make([]byte, n))
func (r *Radix) Itoa(num int, n int) string {
	b := make([]byte, n)
	r.FormatInt(num, b)
	return string(b)
}

// Atoi is equivalent to ParseInt(str)
func (r *Radix) Atoi(str string) (int, error) {
	return r.ParseInt(str)
}
