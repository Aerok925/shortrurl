package reducing

import (
	"fmt"
	"hash/crc32"
)

type Reducing struct {
}

func New() *Reducing {
	return &Reducing{}
}

func (r *Reducing) TruncateLine(str string) string {
	hash := crc32.ChecksumIEEE([]byte(str))
	return fmt.Sprintf("%x", hash)
}
