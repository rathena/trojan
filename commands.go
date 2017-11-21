package integration

// A 2-byte block of data representing all communicated operations.
//
// The format of each PacketCommand is SenderToReceiver, followed a
// succinct description of the requested operation of the function name.
//
// When multiple PacketCommands trigger the same operation, they will
// be suffixed by their numeric code.
//
// When known, in-line comments will describe the complete byte format
// using the syntax defined in `docs/packet_struct_notation.txt` by rathena.
//
// Details for each property name and there data type, as opposed to byte size,
// will be defined in the related functions documentation.
type PacketCommand uint16

// Client to Login Server.
const (
	ClientToLoginRequestAccess0064 PacketCommand = 0x0064 // <version>.L <username>.24B <password>.24B <client_type>.B
)

// Login Server to Client.
const (
	LoginToClientSuccess0069 PacketCommand = 0x0069 // <packet_size>.W <auth_one>.L <account_id>.L <auth_one>.L <unknown>.30B <sex>.B { <world_ip>.W <world_port>.W <world_name>.20B <world_users>.W <world_type>.W <world_new>.W }*
	LoginToClientFailed      PacketCommand = 0x006a // <code>.W <string_empty_or_timestamp>.20B
	LoginToClientError       PacketCommand = 0x0081 // <code>.B
)

// Char Server to Login Server.
const (
	CharToLoginRegister PacketCommand = 0x2710 // <username>.24B <password>.24B <unknown>.L <ip_address>.L <port>.W <world_name>.20B <unknown>.W <world_type>.W <world_new>.W
)

// Login Server to Char Server.
const (
	LoginToCharRegistrationReply PacketCommand = 0x2711 // used with uint8 code 0 for success and 3 for failure
)
