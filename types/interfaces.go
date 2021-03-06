package types

// Config is the generic configuration
type Config interface {
}

type PID interface {
	Update(ImuRotations) []Throttle
}

type ESC interface {
	SetThrottles([]Throttle)
	MotorsOn()
	MotorsOff()
}

// Logger is interface for the udpLogger
type UdpLogger interface {
	Send(ImuRotations)
}

// IMU is interface to imu mems
type IMU interface {
	Close()
	GetRotations() (ImuRotations, error)
	ResetReadingTimes()
	CanRead() bool
}
