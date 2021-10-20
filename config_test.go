package ngamux

import (
	"testing"

	"github.com/golang-must/must"
)

func TestBuildConfig(t *testing.T) {
	must := must.New(t)
	result := buildConfig()
	expected := Config{
		RemoveTrailingSlash: true,
	}
	must.Equal(expected.RemoveTrailingSlash, result.RemoveTrailingSlash)

	result = buildConfig(Config{
		RemoveTrailingSlash: false,
	})
	expected = Config{
		RemoveTrailingSlash: false,
	}
	must.Equal(expected.RemoveTrailingSlash, result.RemoveTrailingSlash)
}
