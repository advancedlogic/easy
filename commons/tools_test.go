package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadLinesOfFile(t *testing.T) {
	lines, err := ReadLinesOfFile("../assets/multilines.txt")
	assert.Equal(t, err, nil)
	assert.Equal(t, len(lines), 3)
}

func TestEB64(t *testing.T) {
	str := EB64("test")
	assert.NotEqual(t, str, "test")
	assert.Equal(t, str, "dGVzdA==---")
}

func TestDB64(t *testing.T) {
	str := DB64("dGVzdA==---")
	assert.NotEqual(t, str, "dGVzdA==---")
	assert.Equal(t, str, "test")
}

func TestUUID(t *testing.T) {
	u := UUID()
	assert.NotEqual(t, u, "")
}

func TestSHA1(t *testing.T) {
	s := SHA1("test")
	assert.NotEqual(t, s, "test")
	assert.Equal(t, s, "qUqP5cyxm6YcTAhz05Hph5gvu9M=")
}
