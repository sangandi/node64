//go:build !linux

package localnet

// Nicht-Linux: Netztyp nicht bekannt; bisherige konservative Bridge-Erkennung behalten.
func interfaceKinds() map[string]string { return nil }
