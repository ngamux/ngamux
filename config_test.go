package ngamux

import (
	"testing"
)

func TestBuildConfig(t *testing.T) {
	result := buildConfig()
	expected := Config{
		RemoveTrailingSlash: true,
	}

	if result.RemoveTrailingSlash != expected.RemoveTrailingSlash {
		t.Errorf("TestBuildConfig need %v, but got %v", expected.RemoveTrailingSlash, result.RemoveTrailingSlash)
	}

	result = buildConfig(Config{
		RemoveTrailingSlash: false,
	})
	expected = Config{
		RemoveTrailingSlash: false,
	}

	if result.RemoveTrailingSlash != expected.RemoveTrailingSlash {
		t.Errorf("TestBuildConfig need %v, but got %v", expected.RemoveTrailingSlash, result.RemoveTrailingSlash)
	}
}
