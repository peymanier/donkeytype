package utils

import "math/rand/v2"

func Sample[T any](src []T, count int) []T {
	dst := make([]T, len(src))
	copy(dst, src)

	rand.Shuffle(len(dst), func(i, j int) {
		dst[i], dst[j] = dst[j], dst[i]
	})

	return dst[:count]
}
