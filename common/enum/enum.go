package enum

const (
	Socks5Version = uint8(5)
)

const (
	ConnectCommand   = uint8(1)
	BindCommand      = uint8(2)
	AssociateCommand = uint8(3)
	Ipv4Address      = uint8(1)
	FqdnAddress      = uint8(3)
	Ipv6Address      = uint8(4)
)

const (
	SuccessReply uint8 = iota
	ServerFailure
	RuleFailure
	NetworkUnreachable
	HostUnreachable
	ConnectionRefused
	TTLExpired
	CommandNotSupported
	AddrTypeNotSupported
)

const (
	NoAuth          = uint8(0)
	NoAcceptable    = uint8(255)
	UserPassAuth    = uint8(2)
	UserAuthVersion = uint8(1)
	AuthSuccess     = uint8(0)
	AuthFailure     = uint8(1)
)

const (
	Timeout    int = 60
	BufferSize int = 4 * 1024
)
