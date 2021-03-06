package v124

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type (
	GcomoRaw struct {
		Agc         uint16 // 0 = manual, 1 = agc on
		Attenuation uint16 // 0-63 (0.5 dB steps)
		Power       uint16 // 0-65535
		TargetPower uint16 // 0-65535 target for AGC
	}

	DownconverterRaw struct {
		Enabled     uint16 // 1 == output off, 2 = output on
		Locked      uint16 // 0 == no lock, 1 = lock
		Attenuation uint16 // 0-31 in dB
		Frequency   uint16 // MHz
	}

	BitStatistics32Raw struct {
		Pattern [4]uint32 // 00, 01, 10, 11
	}

	BitStatistics16Raw struct {
		Pattern [4]uint16 // 00, 01, 10, 11
	}

	Adb3lRaw struct {
		TotalPower       [4]uint32
		BitStatistics    [4]BitStatistics32Raw
		DelayCorrelation [3]int32 // S0-S1, S1-S2, S2-S3, TODO: is this type right
	}

	Core3hRaw struct {
		_                uint32 // VDIF timestamp, unused in version
		PpsDelay         uint32 // in ns
		TotalPowerCalOn  uint32
		TotalPowerCalOff uint32
		_                uint32 //unused
		_                uint32 //unused
	}

	BbcRaw struct {
		Frequency           uint32
		Bandwidth           uint8
		Agc                 uint8
		GainUSB             uint8
		GainLSB             uint8
		TotalPowerUSBCalOn  uint32
		TotalPowerLSBCalOn  uint32
		TotalPowerUSBCalOff uint32
		TotalPowerLSBCalOff uint32
		_                   [4]uint16
		_                   uint16
		_                   uint16
		_                   uint16
		_                   uint16
	}

	Dbbc3DdcMulticastRaw struct {
		Version       [32]byte
		Gcomo         [8]GcomoRaw
		Downconverter [8]DownconverterRaw
		Adb3l         [8]Adb3lRaw
		Core3h        [8]Core3hRaw
		Bbc           [128]BbcRaw
	}
)

func (g GcomoRaw) Cook() Gcomo {
	return Gcomo{
		Agc:         g.Agc == 1,
		Attenuation: float64(g.Attenuation) / 2,
		Power:       g.Power,
		TargetPower: g.TargetPower,
	}
}

func (d DownconverterRaw) Cook() Downconverter {
	return Downconverter{
		Enabled:     d.Enabled == 2,
		Locked:      d.Locked == 1,
		Attenuation: float64(d.Attenuation),
		Frequency:   float64(d.Frequency),
	}
}

func (a Adb3lRaw) Cook() Adb3l {
	bits := make([]map[string]uint32, len(a.BitStatistics))

	for i, bs := range a.BitStatistics {
		bits[i] = make(map[string]uint32)
		for j := range bs.Pattern {
			bits[i][fmt.Sprintf("%02b", j)] = bs.Pattern[j]
		}

	}
	return Adb3l{
		TotalPower:       a.TotalPower[:],
		DelayCorrelation: a.DelayCorrelation[:],
		BitStatistics:    bits,
	}
}

func (c Core3hRaw) Cook() Core3h {
	return Core3h{
		PpsDelay:         c.PpsDelay,
		TotalPowerCalOn:  c.TotalPowerCalOn,
		TotalPowerCalOff: c.TotalPowerCalOff,
	}
}

func (b BbcRaw) Cook() Bbc {
	return Bbc{
		Frequency:           float64(b.Frequency) / 524288,
		Bandwidth:           b.Bandwidth,
		Agc:                 b.Agc,
		GainUSB:             b.GainUSB,
		GainLSB:             b.GainLSB,
		TotalPowerUSBCalOn:  b.TotalPowerUSBCalOn,
		TotalPowerLSBCalOn:  b.TotalPowerLSBCalOn,
		TotalPowerUSBCalOff: b.TotalPowerUSBCalOff,
		TotalPowerLSBCalOff: b.TotalPowerLSBCalOff,
	}
}

type (
	Gcomo struct {
		Agc         bool    `json:"agc"`          // 0 = manual, 1 = agc on
		Attenuation float64 `json:"attenuation"`  // 0-63 (0.5 dB steps)
		Power       uint16  `json:"power"`        // 0-65535
		TargetPower uint16  `json:"target_power"` // 0-65535 target for AGC
	}

	Downconverter struct {
		Enabled     bool    `json:"enabled"`     // 1 == output off, 2 = output on
		Locked      bool    `json:"locked"`      // 0 == no lock, 1 = lock
		Attenuation float64 `json:"attenuation"` // 0-31 in dB
		Frequency   float64 `json:"frequency"`   // MHz
	}

	Adb3l struct {
		TotalPower       []uint32            `json:"total_power"`
		BitStatistics    []map[string]uint32 `json:"bit_statistics"`
		DelayCorrelation []int32             `json:"delay_correlation"`
	}

	Core3h struct {
		Timestamp        uint32 `json:"timestamp"` // VDIF timestamp
		PpsDelay         uint32 `json:"pps_delay"` // in ns
		TotalPowerCalOn  uint32 `json:"total_power_cal_on"`
		TotalPowerCalOff uint32 `json:"total_power_cal_off"`
	}

	Bbc struct {
		Frequency           float64 `json:"frequency"`
		Bandwidth           uint8   `json:"bandwidth"`
		Agc                 uint8   `json:"agc"`
		GainUSB             uint8   `json:"gain_usb"`
		GainLSB             uint8   `json:"gain_lsb"`
		TotalPowerUSBCalOn  uint32  `json:"total_power_usb_cal_on"`
		TotalPowerLSBCalOn  uint32  `json:"total_power_lsb_cal_on"`
		TotalPowerUSBCalOff uint32  `json:"total_power_usb_cal_off"`
		TotalPowerLSBCalOff uint32  `json:"total_power_lsb_cal_off"`
	}

	Dbbc3DdcMulticast struct {
		Version       string          `json:"version"`
		Gcomo         []Gcomo         `json:"gcomo"`
		Downconverter []Downconverter `json:"downconverter"`
		Adb3l         []Adb3l         `json:"adb3l"`
		Core3h        []Core3h        `json:"core3h"`
		Bbc           []Bbc           `json:"bbc"`
	}
)

func (d Dbbc3DdcMulticastRaw) Cook() Dbbc3DdcMulticast {
	// TODO: handle number of boards
	dc := Dbbc3DdcMulticast{}

	dc.Version = cstr(d.Version[:])

	dc.Gcomo = make([]Gcomo, len(d.Gcomo))
	for i, g := range d.Gcomo {
		dc.Gcomo[i] = g.Cook()
	}

	dc.Downconverter = make([]Downconverter, len(d.Downconverter))
	for i, d := range d.Downconverter {
		dc.Downconverter[i] = d.Cook()
	}

	dc.Adb3l = make([]Adb3l, len(d.Adb3l))
	for i, d := range d.Adb3l {
		dc.Adb3l[i] = d.Cook()
	}

	dc.Core3h = make([]Core3h, len(d.Core3h))
	for i, d := range d.Core3h {
		dc.Core3h[i] = d.Cook()
	}

	dc.Bbc = make([]Bbc, len(d.Bbc))
	for i, d := range d.Bbc {
		dc.Bbc[i] = d.Cook()
	}

	return dc
}

func cstr(str []byte) string {
	for n, b := range str {
		if b == 0 {
			return string(str[:n])
		}
	}
	return string(str)
}

func (d *Dbbc3DdcMulticast) UnmarshalBinary(buf []byte) error {
	var dRaw Dbbc3DdcMulticastRaw

	reader := bytes.NewReader(buf)
	err := binary.Read(reader, binary.LittleEndian, &dRaw)
	if err != nil {
		return err
	}
	*d = dRaw.Cook()

	return nil
}
