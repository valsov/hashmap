package hashmap

import (
	"encoding/binary"
	"unsafe"
)

const RAPID_SEED = 0xbdd89aa982704029

var rapidSecret = [3]uint64{0x2d358dccaa6c78a5, 0x8bb84b93962eacc9, 0x4b33a62ed433d4a3}

func Hash(p unsafe.Pointer, length int) uint64 {
	// Get key bytes
	key := unsafe.Slice((*byte)(p), length)
	// Hash key
	return rapidhashInternal(key, length, RAPID_SEED, rapidSecret)
}

func rapidhashInternal(key []byte, length int, seed uint64, secret [3]uint64) uint64 {
	seed ^= rapidMix(seed^secret[0], secret[1]) ^ uint64(length)
	var a, b uint64

	if length <= 16 {
		if length >= 4 {
			plast := length - 4
			a = (rapidRead32(key) << 32) | rapidRead32(key[plast:])
			delta := (length & 24) >> (length >> 3)
			b = (rapidRead32(key[delta:]) << 32) | rapidRead32(key[plast-delta:])
		} else if length > 0 {
			a = rapidReadSmall(key, length)
			b = 0
		} else {
			a = 0
			b = 0
		}
	} else {
		i := length
		if i > 48 {
			see1, see2 := seed, seed
			for i >= 96 {
				seed = rapidMix(rapidRead64(key)^secret[0], rapidRead64(key[8:])^seed)
				see1 = rapidMix(rapidRead64(key[16:])^secret[1], rapidRead64(key[24:])^see1)
				see2 = rapidMix(rapidRead64(key[32:])^secret[2], rapidRead64(key[40:])^see2)
				seed = rapidMix(rapidRead64(key[48:])^secret[0], rapidRead64(key[56:])^seed)
				see1 = rapidMix(rapidRead64(key[64:])^secret[1], rapidRead64(key[72:])^see1)
				see2 = rapidMix(rapidRead64(key[80:])^secret[2], rapidRead64(key[88:])^see2)
				key = key[96:]
				i -= 96
			}
			if i >= 48 {
				seed = rapidMix(rapidRead64(key)^secret[0], rapidRead64(key[8:])^seed)
				see1 = rapidMix(rapidRead64(key[16:])^secret[1], rapidRead64(key[24:])^see1)
				see2 = rapidMix(rapidRead64(key[32:])^secret[2], rapidRead64(key[40:])^see2)
				key = key[48:]
				i -= 48
			}
			seed ^= see1 ^ see2
		}
		if i > 16 {
			seed = rapidMix(rapidRead64(key)^secret[2], rapidRead64(key[8:])^seed^secret[1])
			if i > 32 {
				seed = rapidMix(rapidRead64(key[16:])^secret[2], rapidRead64(key[24:])^seed)
			}
		}
		a = rapidRead64(key[i-16:])
		b = rapidRead64(key[i-8:])
	}
	a ^= secret[1]
	b ^= seed
	rapidMum(&a, &b)
	return rapidMix(a^secret[0]^uint64(length), b^secret[1])
}

func rapidMum(A, B *uint64) {
	ha, hb := *A>>32, *B>>32
	la, lb := uint32(*A), uint32(*B)
	rh := ha * hb
	rm0 := ha * uint64(lb)
	rm1 := hb * uint64(la)
	rl := uint64(la) * uint64(lb)
	t := rl + (rm0 << 32)
	c := uint64(0)
	if t < rl {
		c = 1
	}
	lo := t + (rm1 << 32)
	if lo < t {
		c++
	}
	hi := rh + (rm0 >> 32) + (rm1 >> 32) + c
	*A = lo
	*B = hi
}

func rapidMix(A, B uint64) uint64 {
	rapidMum(&A, &B)
	return A ^ B
}

func rapidRead64(p []byte) uint64 {
	return binary.LittleEndian.Uint64(p)
}

func rapidRead32(p []byte) uint64 {
	return uint64(binary.LittleEndian.Uint32(p))
}

func rapidReadSmall(p []byte, k int) uint64 {
	return (uint64(p[0]) << 56) | (uint64(p[k>>1]) << 32) | uint64(p[k-1])
}

/*
//
// rapidSecret computing
//
func makeSecret(seed uint64, secret *[4]uint64) {
	c := []uint8{15, 23, 27, 29, 30, 39, 43, 45, 46, 51, 53, 54, 57, 58, 60, 71, 75, 77, 78, 83, 85, 86, 89, 90, 92, 99, 101, 102, 105, 106, 108, 113, 114, 116, 120, 135, 139, 141, 142, 147, 149, 150, 153, 154, 156, 163, 165, 166, 169, 170, 172, 177, 178, 180, 184, 195, 197, 198, 201, 202, 204, 209, 210, 212, 216, 225, 226, 228, 232, 240}

	for i := 0; i < 4; i++ {
		ok := false
		for !ok {
			ok = true
			secret[i] = 0
			for j := 0; j < 64; j += 8 {
				secret[i] |= uint64(c[wyrand(&seed)%uint64(len(c))]) << j
			}
			if secret[i]%2 == 0 {
				ok = false
				continue
			}
			for j := 0; j < i; j++ {
				if bits.OnesCount64(secret[j]^secret[i]) != 32 {
					ok = false
					break
				}
			}
			if ok && !isPrime(secret[i]) {
				ok = false
			}
		}
	}
}

func wyrand(seed *uint64) uint64 {
	*seed += 0x2d358dccaa6c78a5
	return rapidMix(*seed, *seed^0x8bb84b93962eacc9)
}

func isPrime(n uint64) bool {
	if n < 2 || n&1 == 0 {
		return false
	}
	if n < 4 {
		return true
	}
	if !sprp(n, 2) {
		return false
	}
	if n < 2047 {
		return true
	}
	for _, a := range []uint64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37} {
		if !sprp(n, a) {
			return false
		}
	}
	return true
}

func mulMod(a, b, m uint64) uint64 {
	var r uint64
	for b > 0 {
		if b&1 != 0 {
			r2 := r + a
			if r2 < r {
				r2 -= m
			}
			r = r2 % m
		}
		b >>= 1
		if b > 0 {
			a2 := a + a
			if a2 < a {
				a2 -= m
			}
			a = a2 % m
		}
	}
	return r
}

func powMod(a, b, m uint64) uint64 {
	r := uint64(1)
	for b > 0 {
		if b&1 != 0 {
			r = mulMod(r, a, m)
		}
		b >>= 1
		if b > 0 {
			a = mulMod(a, a, m)
		}
	}
	return r
}

func sprp(n, a uint64) bool {
	d := n - 1
	s := uint(0)
	for d&0xff == 0 {
		d >>= 8
		s += 8
	}
	if d&0xf == 0 {
		d >>= 4
		s += 4
	}
	if d&0x3 == 0 {
		d >>= 2
		s += 2
	}
	if d&0x1 == 0 {
		d >>= 1
		s += 1
	}
	b := powMod(a, d, n)
	if b == 1 || b == n-1 {
		return true
	}
	for r := uint(1); r < s; r++ {
		b = mulMod(b, b, n)
		if b <= 1 {
			return false
		}
		if b == n-1 {
			return true
		}
	}
	return false
}

*/
