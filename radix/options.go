package radix

var (
	charset = []byte{
		'A', 'B', 'E', 'F', 'r', 's', 't', 'H',
		'm', 'n', 'p', 'q', 'G', 'u', 'P', 'D',
		'J', 'K', 'M', 'N', 'R', 'g', 'C', 'S',
		'c', 'd', 'e', 'f', 'Q', 'h', 'j', 'k',
		'2', '3', '4', '5', 'T', 'U', '9', 'W',
		'8', 'Y', 'a', 'b', '6', 'w', '7', 'y',
	}

	offset = func(idx int) int { return idx * 7 }
)

type apply func(r *Radix)

func newOptions(opts ...apply) []apply {
	return append([]apply{
		Charset(charset), Offset(offset),
	}, opts...)
}

func Charset(charset []byte) apply {
	return func(r *Radix) {
		charsetRef := make(map[byte]int)
		for i, b := range charset {
			charsetRef[b] = i
		}
		radix, radixMod := len(charset), 0
		if radix&(radix-1) == 0 {
			radixMod = radix - 1
		}
		r.charset = charset
		r.charsetRef = charsetRef
		r.radix = radix
		r.radixMod = radixMod
	}
}

func Offset(offset func(i int) int) apply {
	return func(r *Radix) {
		r.offset = offset
	}
}
