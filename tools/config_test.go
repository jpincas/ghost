// Copyright 2017 Jonathan Pincas

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ghost

import "testing"

func TestCompareBundles(t *testing.T) {
	cases := []struct {
		bundles1, bundles2 Bundles
		want               bool
	}{
		{Bundles([]string{"bundle1"}), Bundles([]string{"bundle1"}), true},
		{Bundles([]string{}), Bundles([]string{}), true},
		{Bundles([]string{"bundle1", "bundle2", "bundle3"}), Bundles([]string{"bundle1", "bundle3"}), false},
		{Bundles([]string{"bundle1", "bundle3"}), Bundles([]string{"bundle1", "bundle2"}), false},
	}

	for _, c := range cases {
		got := compareBundles(c.bundles1, c.bundles2)
		if got != c.want {
			t.Fatalf("compareBundles(%q, %q) == %v, want %v", c.bundles1, c.bundles2, got, c.want)
		}
	}
}

func TestInstallBundle(t *testing.T) {

	cases := []struct {
		oldBundles  Bundles
		newBundle   string
		wantBundles Bundles
	}{
		{*new(Bundles), "bundle1", Bundles([]string{"bundle1"})},
		{Bundles([]string{"bundle1"}), "bundle2", Bundles([]string{"bundle1", "bundle2"})},
		{Bundles([]string{"bundle1"}), "bundle1", Bundles([]string{"bundle1"})},
	}

	for _, c := range cases {
		gotBundles, _ := c.oldBundles.InstallBundle(c.newBundle)
		if !compareBundles(gotBundles, c.wantBundles) {
			t.Errorf("installBundle(%q) into %q == %q, want %q", c.newBundle, c.oldBundles, gotBundles, c.wantBundles)
		}
	}

}

func TestUninstallBundle(t *testing.T) {

	cases := []struct {
		oldBundles  Bundles
		newBundle   string
		wantBundles Bundles
	}{
		{*new(Bundles), "bundle1", Bundles([]string{})},                                                          //No bundles installed
		{Bundles([]string{"bundle1"}), "bundle2", Bundles([]string{"bundle1"})},                                  //Bundle to be uninstalled not installed
		{Bundles([]string{"bundle1"}), "bundle1", Bundles([]string{})},                                           //Should uninstall the bundle
		{Bundles([]string{"bundle1", "bundle2", "bundle3"}), "bundle2", Bundles([]string{"bundle1", "bundle3"})}, //Should uninstall the bundle
	}

	for _, c := range cases {
		gotBundles, _ := c.oldBundles.UnInstallBundle(c.newBundle)
		if !compareBundles(gotBundles, c.wantBundles) {
			t.Errorf("unInstallBundle(%q) from %q == %q, want %q", c.newBundle, c.oldBundles, gotBundles, c.wantBundles)
		}
	}

}
