package ngamux

import (
	"reflect"
	"testing"
)

func TestNewNgamux(t *testing.T) {
	result := NewNgamux()
	expected := &Ngamux{
		routes:            routeMap{},
		routesParam:       routeMap{},
		config:            buildConfig(),
		regexpParamFinded: paramsFinder,
	}

	if !reflect.DeepEqual(result.routes, expected.routes) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routes, result.routes)
	}

	if !reflect.DeepEqual(result.routesParam, expected.routesParam) {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.routesParam, result.routesParam)
	}

	if result.config.RemoveTrailingSlash != expected.config.RemoveTrailingSlash {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.config.RemoveTrailingSlash, result.config.RemoveTrailingSlash)
	}

	if result.regexpParamFinded != expected.regexpParamFinded {
		t.Errorf("TestNewNgamux need %v, but got %v", expected.regexpParamFinded, result.regexpParamFinded)
	}
}
