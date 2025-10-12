package mod_net

const (
	MaxIPv4Bits    = 32
	MaxIPv6Bits    = 128
	HostSubnetBits = 2 // /30
	UserSubnetBits = 5 // /30
	MaxVLAN        = 4096
	MaxTI          = 16384 // Juniper's JunOS `st0` interface max `unit` number
)
