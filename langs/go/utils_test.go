package gogen

import (
	"testing"
)

func TestGenNamespace(t *testing.T) {
	cases := []struct {
		namespace string
		path      string
		pkgName   string
	}{
		{
			"github.com..ezbuy..tgen..thriftgotest..unusedInclude",
			"github.com/ezbuy/tgen/thriftgotest/unusedInclude",
			"unusedInclude",
		},
		{
			"github.com.ezbuy.tgen.thriftgotest.unusedInclude",
			"github/com/ezbuy/tgen/thriftgotest/unusedInclude",
			"unusedInclude",
		},
	}

	utils := &TplUtils{}

	for _, one := range cases {
		pkgName, path := utils.GenNamespace(one.namespace)

		if path != one.path {
			t.Errorf("expected path: %s, got %s", one.path, path)
		}

		if pkgName != one.pkgName {
			t.Errorf("expected package name: %s, got %s", one.pkgName, pkgName)
		}
	}
}
