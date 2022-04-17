package backend

import (
	"testing"
)

func BenchmarkHashPassword(b *testing.B) {
	salt := generatePasswordSalt()
	for i := 0; i < b.N; i++ {
		hashPassword("OoTec5vi ", salt)
	}
}
