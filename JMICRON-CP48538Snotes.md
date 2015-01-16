The JMICRON-CP48 538S is one of many different USB-SATA bridge chips used by WD MyBooks with encryption. (According to the database used by Western Digital's firmware upgrade tool's file database, there are several models that use this chip AND the firmware described below, but not the encryption...) There appears to be absolutely no mention of this chip on the Internet, let alone documentation.

I won't release the dump that I got from the same place Western Digital's own firmware updater software gets it from, nor will I release the disassembly — at least not on github. Here's what you need to know about the file, though:

```
filename: Release-VS-1025-20130711.bin
size:     49152 bytes
crc32:    cfd13030
md5:      7f75e5d59cfac57579effe2d9d5388de
sha1:     a14c6fb97cfa76e5bd1d007897bde20be132c4b8
```

Due to the complete lack of documentation, we will have to work entirely from scratch here.

The firmware is powered by what appears to be a regular old Intel 8501 core.

IDA doesn't want to play nice with this code so let's pick out important items.

The function at 0x1BF0 takes the four bytes after it on the caller's instruction stream as parameters...

RAM 0x3206 appears to be the start of the key sector...?

The code at 0x613 indicates that our ROM is not a boot ROM and that it should be mapped to code memory at 0x4000. Therefore, in the following, if I prefix an address with ~, it is a virtual address.

TODO does this have combined ROM/RAM? 0x1BF0 is still the function above...

~0x5287 writes the `WDv1` signatures and the two bytes after the first `WDv1` to 0x3206; this is called at ~0x7335, which is subsequently followed by copying those first 0x60 bytes of the key sector to RAM 0x3800...

At ~0xB432 is an array of structures describing devices, of the following form:
```go
type deviceDescriptorValues struct {
	VendorID         uint16
	ProductID        uint16
	unknown          [3]byte
	unkStringPtrs    [3]uint16       // [0] appears to be vendor; [2] appears to be product
}
```
The last element is at ~0xB584. I do not know where this is accessed, or how the firmware knows which device it is being used with currently. (Hardcoded by the firmware updater? This won't be important for our purposes, but finding where this array is used **is**)

I do not see either a CBW or a CSW signature in here, at least not directly...

The function at ~0x5217 sets `DPTR` to whatever's at RAM 0x409C and then checks `DPTR` and `DPTR+0x5C` to see if they match `WDv1`. It then checks a few other things that I don't know about... (checksum of just the first 0x58 bytes?)

## Known boot ROM routines
I believe these are provided by the boot ROM; if they are actually in RAM and copied on system startup, I do not know (TODO).
- 0x1BBC - compare R0 to R4, R1 to R5, R2 to R6, R3 to R7
- 0x2B50 - memcpy
	- R6:R7 - (TODO source or destination? initial observation suggested source but after seeing random register assignment order I'm not sure)
	- R4:R5 - (TODO opposite of above)
	- R3 - length (bytes)
