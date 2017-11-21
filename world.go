package integration

// A World with all dependent properties.
//
// Type effects the server appearance in the client.
//
//	- 0: appears as normal
//	- 1: replaces user count with "(on the maintenance)"
//	- 2: adds "- Over the age of 18" after user count
//	- 3: adds "- Pay to Play" after user count
//	- 4: adds "- Free Server" after user count
//
// Setting `New` to 1 will cause the client to literally prefix the
// server name with `[New Server]`.
type World struct {
	IP    uint32
	Port  uint16
	Name  [20]byte
	Users uint16
	Type  uint16
	New   uint16
}
