package v125

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
		Timestamp        uint32 // VDIF timestamp
		PpsDelay         uint32 // in ns
		TotalPowerCalOn  uint32
		TotalPowerCalOff uint32
		Tsys             uint32
		Sefd             uint32
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
		BitStatistics       BitStatistics16Raw
		TsysUSB             uint16
		TsysLSB             uint16
		SEFDUSB             uint16
		SEFDLSB             uint16
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
		for j := range bs.Pattern {
			bits[i] = make(map[string]uint32)
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
	return Core3h(c)
}

func (b BbcRaw) Cook() Bbc {
	bits := make(map[string]uint32)

	for j := range b.BitStatistics.Pattern {
		bits[fmt.Sprintf("%02b", j)] = uint32(b.BitStatistics.Pattern[j])
	}
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
		BitStatistics:       bits,
		TsysUSB:             b.TsysUSB,
		TsysLSB:             b.TsysLSB,
		SEFDUSB:             b.SEFDUSB,
		SEFDLSB:             b.SEFDLSB,
	}
}

type (
	Gcomo struct {
		Agc         bool    // 0 = manual, 1 = agc on
		Attenuation float64 // 0-63 (0.5 dB steps)
		Power       uint16  // 0-65535
		TargetPower uint16  // 0-65535 target for AGC
	}

	Downconverter struct {
		Enabled     bool    // 1 == output off, 2 = output on
		Locked      bool    // 0 == no lock, 1 = lock
		Attenuation float64 // 0-31 in dB
		Frequency   float64 // MHz
	}

	Adb3l struct {
		TotalPower       []uint32
		BitStatistics    []map[string]uint32
		DelayCorrelation []int32
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
		Frequency           float64
		Bandwidth           uint8
		Agc                 uint8
		GainUSB             uint8
		GainLSB             uint8
		TotalPowerUSBCalOn  uint32
		TotalPowerLSBCalOn  uint32
		TotalPowerUSBCalOff uint32
		TotalPowerLSBCalOff uint32
		BitStatistics       map[string]uint32
		TsysUSB             uint16
		TsysLSB             uint16
		SEFDUSB             uint16
		SEFDLSB             uint16
	}

	Dbbc3DdcMulticast struct {
		Version       string
		Gcomo         []Gcomo
		Downconverter []Downconverter
		Adb3l         []Adb3l
		Core3h        []Core3h
		Bbc           []Bbc
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
