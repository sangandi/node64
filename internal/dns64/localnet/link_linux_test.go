//go:build linux

package localnet

import "testing"

func TestKindFromLinkInfo(t *testing.T) {
	// NLA type 2 (unbekannt), gefolgt von NLA type 1 (IFLA_INFO_KIND="ipvlan").
	data := []byte{8, 0, 2, 0, 'x', 0, 0, 0, 11, 0, 1, 0, 'i', 'p', 'v', 'l', 'a', 'n', 0, 0}
	if got := kindFromLinkInfo(data); got != "ipvlan" {
		t.Fatalf("got %q, want ipvlan", got)
	}
	if got := kindFromLinkInfo([]byte{3, 0, 1, 0}); got != "" {
		t.Fatalf("malformed attribute: got %q", got)
	}
}
