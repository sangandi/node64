package localnet

import "testing"

func TestInBridge(t *testing.T) {
	veth := iface{"12", "13"}
	real := iface{"2", "2"}
	for _, c := range []struct {
		name      string
		container bool
		ifs       map[string]iface
		want      bool
	}{
		{"kein Container", false, map[string]iface{"lo": {"1", "1"}, "eth0": veth}, false},
		{"Bridge", true, map[string]iface{"lo": {"1", "1"}, "eth0": veth}, true},
		{"Host-Netz (echte Karte)", true, map[string]iface{"lo": {"1", "1"}, "enp3s0": real, "veth1": veth}, false},
		{"Host-Netz (docker0)", true, map[string]iface{"lo": {"1", "1"}, "docker0": veth}, false},
		{"nur lo", true, map[string]iface{"lo": {"1", "1"}}, false},
	} {
		if got := inBridge(c.container, c.ifs); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}

// Docker-Bridge muss weiterhin eine Port-Veroeffentlichung verlangen, IPvlan und Macvlan dagegen nicht.
func TestInBridgeWithKinds(t *testing.T) {
	ifs := map[string]iface{"lo": {"1", "1"}, "eth0": {"12", "13"}}
	for _, tc := range []struct {
		name string
		kind string
		want bool
	}{
		{"veth", "veth", true},
		{"ipvlan", "ipvlan", false},
		{"macvlan", "macvlan", false},
		{"unknown", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := inBridgeWithKinds(true, ifs, map[string]string{"eth0": tc.kind}); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
