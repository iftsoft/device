package driver

type DeviceDriver interface {
	InitDevice() error
	StartDevice() error
	DeviceTimer() error
	StopDevice() error
	ClearDevice() error
}
