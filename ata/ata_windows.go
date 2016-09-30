// 30 september 2016
package ata

import (
	"regexp"
	"bytes"
	"encoding/binary"

	"golang.org/x/sys/windows"
)

func sysDrivespecUsage() string {
	return "Either PhysicalDriveNNN where NNN is between 0 and 255 or D: where D is a drive letter (uppercase or lowercase)."
}

// TODO is PhysicalDrive case insensitive?
const rPhysicalDrive = "PhysicalDrive"

// TODO are leading 0s allowed?
var validDrivespecs = []*regexp.Regexp{
	regexp.MustCompile("^" + rPhysicalDrive + "[0-9]$"),
	regexp.MustCompile("^" + rPhysicalDrive + "[1-9][0-9]$"),
	regexp.MustCompile("^" + rPhysicalDrive + "1[0-9][0-9]$"),
	regexp.MustCompile("^" + rPhysicalDrive + "2[0-4][0-9]$"),
	regexp.MustCompile("^" + rPhysicalDrive + "25[0-5]$"),
	regexp.MustCompile("^[A-Za-z]:$"),
}

func validDrivespec(drivespec string) bool {
	for _, r := range validDrivespecs {
		if r.MatchString(drivespec) {
			return true
		}
	}
	return false
}

type sysATA struct {
	handle	windows.Handle
}

func sysOpen(drivespec string) (*ATA, error) {
	if !validDrivespec(drivespec) {
		return nil, InvalidDrivespecError(drivespec)
	}
	drivespec = `\\.\` + drivespec
	wdrivespec, err := UTF16PtrFromString(drivespec)
	if err != nil {
		return nil, err
	}
	a := new(sysATA)
	a.handle, err = windows.CreateFile(wdrivespec,
		windows.GENERIC_READ | windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0, nil)
	if err != nil {
		return nil, err
	}
	return &ATA{a}, nil
}

func (a *sysATA) Close() error {
	return windows.CloseHandle(a.handle)
}

type _ATA_PASS_THROUGH_EX struct {
	Length			uint16
	AtaFlags			uint16
	PathId			uint8
	TargetId			uint8
	Lun				uint8
	ReservedAsUchar	uint8
	DataTransferLength	uint32
	TimeOutValue		uint32
	ReservedAsUlong	uint32
	DataBufferOffset	_ULONG_PTR
	PreviousTaskFile	[8]uint8
	CurrentTaskFile	[8]uint8
}

const (
	_ATA_FLAGS_DRDY_REQUIRED = 0x01
	_ATA_FLAGS_DATA_IN =0x02
	_ATA_FLAGS_DATA_OUT = 0x04
	_ATA_FLAGS_48BIT_COMMAND = 0x08

	_IOCTL_ATA_PASS_THROUGH = 0x04D02C
)

var sizeofPTE = binary.Sizeof(_ATA_PASS_THROUGH_EX{})

func (c *Command28) toNT(flags uint16, buf []byte) (send []byte, err error) {
	var pte _ATA_PASS_THROUGH_EX

	pte.Length = uint16(sizeofPTE)
	// TODO _ATA_FLAGS_DRDY_REQUIRED?
	pte.AtaFlags = flags
	pte.DataTransferLength = uint32(len(buf))
	pte.TimeOutValue = ^uint32(0)		// TODO
	pte.DataBufferOffset = _ULONG_PTR(sizeofPTE)
	pte.CurrentTaskFile[0] = c.Features
	pte.CurrentTaskFile[1] = c.Count
	pte.CurrentTaskFile[2] = c.LBALow
	pte.CurrentTaskFile[3] = c.LBAMid
	pte.CurrentTaskFile[4] = c.LBAHigh
	pte.CurrentTaskFile[5] = c.Device
	pte.CurrentTaskFile[6] = c.Command

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, pte)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Command28) toNTRead(outbuf []byte) (send []byte, recv []byte, err error) {
	send, err = c.toNT(_ATA_FLAGS_DATA_IN, outbuf)
	if err != nil {
		return nil, nil, err
	}
	recv = make([]byte, sizeofPTE + len(outbuf))
	return send, recv, nil
}

func (c *Command28) toNTWrite(inbuf []byte) (send []byte, recv []byte, err error) {
	s, err := c.toNT(_ATA_FLAGS_DATA_OUT, inbuf)
	if err != nil {
		return nil, nil, err
	}
	send = make([]byte, len(s) + len(inbuf))
	copy(send[:len(s)], s)
	copy(send[len(s):], inbuf)
	recv = make([]byte, sizeofPTE)
	return send, recv, nil
}

type out struct {
	recv		[]byte
	pte		_ATA_PASS_THROUGH_EX
	resp		*Response28
}

func (a *sysATA) perform(send []byte, recv []byte, err error) (o *out, err error) {
	var count uint32
	var pte _ATA_PASS_THROUGH_EX

	if err != nil {
		return nil, err
	}

	err = windows.DeviceIoControl(a.handle,
		_IOCTL_ATA_PASS_THROUGH,
		&send[0],
		uint32(len(send)),
		&recv[0],
		uint32(len(recv)),
		&count,
		nil)
	if err != nil {
		return nil, err
	}
	recv = recv[:count]

	r := bytes.NewReader(recv[:sizeofPTE])
	err = binary.Read(r, binary.LittleEndian, &pte)
	if err != nil {
		return nil, err
	}

	resp := new(Response28)
	resp.Error = pte.CurrentTaskFile[0]
	resp.Count = pte.CurrentTaskFile[1]
	resp.LBALow = pte.CurrentTaskFile[2]
	resp.LBAMid = pte.CurrentTaskFile[3]
	resp.LBAHigh = pte.CurrentTaskFile[4]
	resp.Device = pte.CurrentTaskFile[5]
	resp.Status = pte.CurrentTaskFile[6]

	return &out{
		recv:		recv[pte.DataBufferOffset:],
		pte:		pte,
		resp:		resp,
	}, nil
}

func (a *sysATA) Read28(c *Command28, b []byte) (resp *Response28, n int, err error) {
	out, err := a.perform(c.toNTRead(b))
	if err != nil {
		return nil, 0, err
	}
	copy(b, out.recv)
	return out.resp, len(out.recv), nil
}

func (a *sysATA) Write28(c *Command28, b []byte) (resp *Response28, err error) {
	out, err := a.perform(c.toNTWrite(b))
	if err != nil {
		return nil, err
	}
	return out.resp, nil
}
