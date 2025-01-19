package ngamux

import (
	"testing"

	"github.com/golang-must/must"
)

func TestBuildConfig(t *testing.T) {
	must := must.New(t)
	result := NewConfig()
	expected := Config{
		RemoveTrailingSlash: true,
	}
	must.Equal(expected.RemoveTrailingSlash, result.RemoveTrailingSlash)
}
