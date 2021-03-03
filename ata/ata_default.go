// 30 september 2016

// +build !windows

package ata

// Systems that will only ever have this version of the package:
// - OS X: https://developer.apple.com/library/content/qa/qa1179/_index.html

// TODO show a not yet implemented error on systems that are not OS X

func sysDrivespecUsage() string {
	return "This OS does not support sending arbitrary ATA commands to a disk; consequently, there are no valid drivespecs and any command requiring a drivespec will fail."
}

type sysATA struct{}

func sysOpen(drivespec string) (*ATA, error) {
	return nil, ErrUnsupportedOS
}

func (a *sysATA) Close() error {
	return ErrUnsupportedOS
}

func (a *sysATA) Read28(c *Command28, b []byte) (resp *Response28, n int, err error) {
	return nil, 0, ErrUnsupportedOS
}

func (a *sysATA) Write28(c *Command28, b []byte) (resp *Response28, err error) {
	return nil, ErrUnsupportedOS
}
