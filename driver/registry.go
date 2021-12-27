package driver

import "errors"

type registrySlot struct {
	spec  *Specification
	maker DrvFactory
}

var driverMap map[string]registrySlot

func RegisterDriver(spec *Specification, maker DrvFactory) error {
	if spec == nil || maker == nil {
		return errors.New("bad params")
	}
	if driverMap == nil {
		driverMap = make(map[string]registrySlot)
	}
	slot := registrySlot{
		spec:  spec,
		maker: maker,
	}
	driverMap[spec.ModelName] = slot

	return nil
}

func GetDriverList() []*Specification {
	list := make([]*Specification, 0)
	for _, slot := range driverMap {
		list = append(list, slot.spec)
	}
	return list
}

func GetDeviceDriver(model string) DeviceDriver {
	if slot, ok := driverMap[model]; ok {
		return slot.maker()
	}
	return nil
}

func GetDriverVersions(model string) []string {
	if slot, ok := driverMap[model]; ok {
		return slot.spec.Versions
	}
	return nil
}
