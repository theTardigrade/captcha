package captcha

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"
)

func randomSeed() (n int64) {
	var seedBytes [8]byte

	if _, err := cryptoRand.Read(seedBytes[:]); err == nil {
		n = int64(binary.LittleEndian.Uint64(seedBytes[:]))
	} else {
		n = time.Now().UTC().UnixNano()
	}

	return
}

func randomNew() *rand.Rand {
	return rand.New(rand.NewSource(randomSeed()))
}
