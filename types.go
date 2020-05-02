package dbbc3mcast

import "encoding/json"

type (
	Gcomo struct {
		Agc         uint16 // 0 = manual, 1 = agc on
		Attenuation uint16 // 0-63 (0.5 dB steps)
		Power       uint16 // 0-65535
		TargetPower uint16 // 0-65535 target for AGC
	}

	Downconverter struct {
		Enabled     uint16 // 1 == output off, 2 = output on
		Lock        uint16 // 0 == no lock, 1 = lock
		Attenuation uint16 // 0-31 in dB
		Frequency   uint16 // MHz
	}

	BitStatistics32 struct {
		Pattern [4]uint32 // 00, 01, 10, 11
	}

	BitStatistics16 struct {
		Pattern [4]uint16 // 00, 01, 10, 11
	}

	Adb3l struct {
		TotalPower       [4]uint32
		BitStatistics    [4]BitStatistics32
		DelayCorrelation [3]int32 // S0-S1, S1-S2, S2-S3, TODO: is this type right
	}

	Core3h struct {
		Timestamp        uint32 // VDIF timestamp
		PpsDelay         uint32 // in ns
		TotalPowerCalOn  uint32
		TotalPowerCalOff uint32
		Tsys             uint32
		Sefd             uint32
	}

	Bbc struct {
		Frequency           uint32
		Bandwidth           uint8
		Agc                 uint8
		GainUSB             uint8
		GainLSB             uint8
		TotalPowerUSBCalOn  uint32
		TotalPowerLSBCalOn  uint32
		TotalPowerUSBCalOff uint32
		TotalPowerLSBCalOff uint32
		BitStatistics       BitStatistics16
		TsysUSB             uint16
		TsysLSB             uint16
		SEFDUSB             uint16
		SEFDLSB             uint16
	}

	Dbbc3DdcMulticast struct {
		Version       [32]byte
		Gcomo         [8]Gcomo
		Downconverter [8]Downconverter
		Adb3l         [8]Adb3l
		Core3h        [8]Core3h
		Bbc           [128]Bbc
	}
)

func cstr(str []byte) string {
	for n, b := range str {
		if b == 0 {
			return string(str[:n])
		}
	}
	return string(str)
}

func (d *Dbbc3DdcMulticast) MarshalJSON() ([]byte, error) {
	// only needed to keep null bytes out of Version
	return json.Marshal(struct {
		Version       string
		Gcomo         [8]Gcomo
		Downconverter [8]Downconverter
		Adb3l         [8]Adb3l
		Core3h        [8]Core3h
		Bbc           [128]Bbc
	}{
		cstr(d.Version[:]),
		d.Gcomo,
		d.Downconverter,
		d.Adb3l,
		d.Core3h,
		d.Bbc,
	})
}
