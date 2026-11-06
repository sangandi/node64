//go:build linux

package localnet

import (
	"bytes"
	"encoding/binary"
	"syscall"
)

// interfaceKinds liest den echten Linux-Link-Typ (veth, ipvlan, macvlan) ueber RTM_GETLINK.
// ifindex/iflink alleine koennen diese virtuellen Schnittstellen nicht unterscheiden.
// Falls der Netlink-Zugriff fehlschlaegt, wird konservativ die Bridge-Pruefung beibehalten.
func interfaceKinds() map[string]string {
	b, err := syscall.NetlinkRIB(syscall.RTM_GETLINK, syscall.AF_UNSPEC)
	if err != nil {
		return nil
	}
	messages, err := syscall.ParseNetlinkMessage(b)
	if err != nil {
		return nil
	}
	out := make(map[string]string)
	for _, msg := range messages {
		if msg.Header.Type != syscall.RTM_NEWLINK {
			continue
		}
		attrs, err := syscall.ParseNetlinkRouteAttr(&msg)
		if err != nil {
			continue
		}
		var name, kind string
		for _, a := range attrs {
			switch a.Attr.Type & 0x3fff { // NLA_TYPE_MASK, ohne NLA_F_NESTED
			case syscall.IFLA_IFNAME:
				name = string(bytes.TrimRight(a.Value, "\x00"))
			case syscall.IFLA_LINKINFO:
				kind = kindFromLinkInfo(a.Value)
			}
		}
		if name != "" {
			out[name] = kind
		}
	}
	return out
}

// kindFromLinkInfo liest IFLA_INFO_KIND (Typ 1) aus verschachtelten Netlink-Attributen.
func kindFromLinkInfo(data []byte) string {
	for len(data) >= 4 {
		length := int(binary.NativeEndian.Uint16(data[:2]))
		typeID := binary.NativeEndian.Uint16(data[2:4]) & 0x3fff
		if length < 4 || length > len(data) {
			return ""
		}
		if typeID == 1 { // IFLA_INFO_KIND
			return string(bytes.TrimRight(data[4:length], "\x00"))
		}
		step := (length + 3) &^ 3 // NLA_ALIGNTO = 4
		if step > len(data) {
			break
		}
		data = data[step:]
	}
	return ""
}
