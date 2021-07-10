package types

// Config is the generic configuration
type Config interface {
}

type PidConfig struct {
	ProportionalGain float32 `yaml:"proportional–gain"`
	IntegralGain     float32 `yaml:"integral-gain"`
	DerivativeGain   float32 `yaml:"derivative-gain"`
}

type EscConfig struct {
	MaxPulseWidth float32 `yaml:"max_esc_pulse_width_ms"`
}

type ImuConfig struct {
	ImuDataPerSecond            int     `yaml:"imu_data_per_second"`
	AccLowPassFilterCoefficient float64 `yaml:"acc_lowpass_filter_coefficient"`
	LowPassFilterCoefficient    float64 `yaml:"lowpass_filter_coefficient"`
}

type FlightConfig struct {
	PID PidConfig `yaml:"pid"`
	Imu ImuConfig `yaml:"imu"`
	Esc EscConfig `yaml:"esc"`
}

type UdpLoggerConfig struct {
	Enabled          bool   `yaml:"enabled"`
	IP               string `yaml:"ip"`
	Port             int    `yaml:"port"`
	PacketsPerSecond int    `yaml:"udp_packets_per_second"`
	DataPerSecond    int    `yaml:"udp_data_per_second"`
}

// XYZ is X, Y, Z data
type XYZ struct {
	X, Y, Z float64
}

type RotationsChanges struct {
	DRoll, DPitch, DYaw float64
}

type SensorData struct {
	Error error
	Data  XYZ
}

type Rotations struct {
	Roll, Pitch, Yaw float64
}

type ImuRotations struct {
	Accelerometer Rotations
	Gyroscope     Rotations
	Rotations     Rotations
	ReadTime      int64
	ReadInterval  int64
}

// Sensor is devices that read data in x, y, z format
type Sensor struct {
	Type   string
	Config Config
}

// CommandParameters is parameters for the command
type CommandParameters interface {
}

type Command struct {
	Command    string
	Parameters CommandParameters
}

// GetConfig reads the config
func (a *Sensor) GetConfig() Config {
	return a.Config
}

// SetConfig sets the config
func (a *Sensor) SetConfig(config Config) {
	a.Config = config
}

// IMU is interface to imu mems
type IMU interface {
	Close()
	GetRotations() (ImuRotations, error)
	ResetReadingTimes()
	CanRead() bool
}

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(ImuRotations)
}

type Throttle struct {
	Motor int
	Value float32
}

type PID interface {
	Update(ImuRotations) []Throttle
}

type ESCsHandler interface {
	SetThrottles([]Throttle)
}
