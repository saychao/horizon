package utf8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrub(t *testing.T) {
	assert.Equal(t, "c78c167633163b96b03d1681911b1f9e1babe6d30d3e537c1c7f50af8e33c49", Scrub("c78c167633163b96b03d1681911b1f9e1babe6d30d3e537c1c7f50af8e33c49\u0000"))
	assert.Equal(t, "scott", Scrub("scott"))
	assert.Equal(t, "scött", Scrub("scött"))
	assert.Equal(t, "�(", Scrub(string([]byte{0xC3, 0x28})))
}
