package domain

import (
	"errors"
	"strings"
)

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	minChar  = '0'
	maxChar  = 'z'
	midChar  = 'm'
)

// RankBetween calculates a lexicographically orderable position string strictly between prev and next.
// If prev is empty, it finds a string before next.
// If next is empty, it finds a string after prev.
// If both are empty, it returns the standard initial position "m".
func RankBetween(prev, next string) (string, error) {
	if prev != "" && next != "" && prev >= next {
		return "", errors.New("invalid rank range: prev must be strictly smaller than next")
	}

	if prev == "" && next == "" {
		return string(midChar), nil
	}

	if prev == "" {
		// Place before next
		for i := 0; i < len(next); i++ {
			char := next[i]
			if char > minChar {
				mid := minChar + (char-minChar)/2
				if mid > minChar {
					return next[:i] + string(mid), nil
				}
				return next[:i] + string(minChar) + string(midChar), nil
			}
		}
		return string(minChar) + string(midChar), nil
	}

	if next == "" {
		// Place after prev
		return prev + string(midChar), nil
	}

	// Between prev and next
	var b strings.Builder
	idx := 0
	for {
		var pChar, nChar byte = minChar, maxChar
		if idx < len(prev) {
			pChar = prev[idx]
		}
		if idx < len(next) {
			nChar = next[idx]
		}

		if pChar == nChar {
			b.WriteByte(pChar)
			idx++
			continue
		}

		if nChar-pChar > 1 {
			mid := pChar + (nChar-pChar)/2
			b.WriteByte(mid)
			return b.String(), nil
		}

		// Difference is 1 (adjacent characters)
		b.WriteByte(pChar)
		idx++
		for {
			pChar = minChar
			if idx < len(prev) {
				pChar = prev[idx]
			}
			if pChar < maxChar {
				mid := pChar + (maxChar-pChar)/2
				if mid > pChar {
					b.WriteByte(mid)
					return b.String(), nil
				}
			}
			b.WriteByte(pChar)
			idx++
		}
	}
}