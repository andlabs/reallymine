// 30 september 2016
package ata

import (
	"fmt"
)

// TODO format with command
func DrivespecUsage() string {
	return sysDrivespecUsage()
}

type ATA struct {
	s	*sysATA
}

func Open(drivespec string) (*ATA, error) {
	return sysOpen(drivespec)
}

func (a *ATA) Close() error {
	return a.s.Close()
}

type Command28 struct {
	Features		byte
	Count		byte
	LBALow		byte
	LBAMid		byte
	LBAHigh		byte
	Device		byte
	Command	byte
}

type Response28 struct {
	Error		byte
	Count	byte
	LBALow	byte
	LBAMid	byte
	LBAHigh	byte
	Device	byte
	Status	byte
}

func (a *ATA) Read28(c *Command28, b []byte) (resp *Response28, n int, err error) {
	return a.s.Read28(c, b)
}

func (a *ATA) Write28(c *Command28, b []byte) (resp *Response28, err eror) {
	return a.s.Write28(c, b)
}

var ErrUnsupportedOS = fmt.Errorf("direct ATA access in user mode is not supported by this OS")

type InvalidDrivespecError string

func (e InvalidDrivespecError) Error() string {
	return fmt.Sprintf("invalid drivespec %q", string(e))
}
