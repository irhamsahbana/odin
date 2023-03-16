package adapters

// NetworkDB contains options common to SQL databases accessed over network.
type NetworkDB struct {
	Database string
	User     string
	Password string
	Host     string
	Port     uint16

	ConnectionTimeout int // Seconds

	MaxOpenCons     int // default: 5
	MaxIdleCons     int // default: 5
	ConnMaxLifetime int // Seconds, default: not set
}
