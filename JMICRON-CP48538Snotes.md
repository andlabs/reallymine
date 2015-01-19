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

The function at ~0x44BB writes something to 0x70B6, then checks 0x70B7 repeatedly until bit 0 is set.

The function at ~0x81F3 converts a SCSI READ(10)-format input block to an SCSI READ(16)-format block.

```
KEY SECTOR STRUCTURE
Offset	Size	What		RAM		Code
CHUNK 0 - HEADER (0x60 BYTES)
0		4		WDv1		0x3206	[TODO]
4		2		checksum	0x320A	[TODO]
...
0x20	1		01			0x3226	[TODO]
...
0x24	4		[00 00]FP	0x322A	[TODO]
...
0x30	1		00			0x3236	[TODO]
0x31	1		see next	0x3237	[TODO]
	bit 0 - ??
	bit 1 - set if bit 5 of 0x3F2F is set
	bit 2-7 - ??
0x32	1		FF			0x3238	[TODO]
0x33	1		00			0x3239	[TODO]
...
0x5C	4		WDv1		0x3262	[TODO]
CHUNK 1 - SOMETHING ELSE (0xXXX(0x30?0x90?) BYTES)
0x60	4		WDq1		0x3266	~0x543B
0x64
0x68	0x10	????		0x326E	[TODO]
0x68	2		checksum	0x326E	[TODO] WTF?
...
0x20	1		00			0x3286	[TODO]
	(other bits TODO; set before bit 5)
	bit 5 - set by ~0x5566 if bit 2 of 0x3F2F set
	bit 6-7 - ????
0x21	1		00			0x3287	~0x556D
...
[CHUNK 2?]
0x90	4		WDqe		0x3296	[TODO]
```

The function at ~0x457F is memset. R6:R7 is the destination, R5 is the byte to fill with, and R3 is the number of bytes.

At ~0xC547 are three sets of suspicious 16-byte blocks...

Checks for `WDq1` at 0x3C00 are at ~0x5164, ~0x5194, ~0x5300, and ~0x5390.

Is ~0xBA7E the firmware's main loop? If so 0x3F44 is the current argument. It might not be...
>Upon further investigation it may actually be the key sector writing loop...?

At ~0x7FD7 is an array of byte-address pairs that seem to map to a list of SCSI commands...? The array seems to end with the ele
ment at either ~0x8046 or ~0x8049...
>Another such list is at ~0x810E, ending at either ~0x8150 or ~0x8153.
>A third at ~0x93B1, ends either at ~0x93C9 or ~0x93CC.
>And a fourth at ~0xAAF4, ending either at ~0xAB39 or ~0xAB3C.

According to that, 0x3E9F is the byte that holds the SCSI command.

~0x81F3 is the first function called by one of the READ(10) handlers (~0x8084); this function seems to handle parameters.

TODO is 0x70B0 the ultimate resting place of the starting LBA? Or is 0x4046? or 0x3ED1?! so much copying...

Let's reverse the path through the code at ~0x7FCE, which seems to contain the only case of a SCSI READ command in the SCSI command tables (the ~0x7FD7 one).
- ~0x7FCE is a subroutine called at ~0x9356
- Immediately prior to that, 0x3E9F, the 16 input block bytes, are copied from 0x300F.
- Immediately prior to that is a check that determines whether any of the above is even run: the function at ~0x97AC is called; if it returns with the carry bit SET, then the above code is SKIPPED. This routine does:
```
if xdata 0x7037 ^ 1 == 0
	return carry clear
if xdata 0x7036 != 0x41 && xdata 0x7036 != 0x81
	return carry clear
if xdata 0x300F & 0xF0 == 0x80		// SCSI command
	irrelevant; now we know data is already read by this point
if xdata 0x300F & 0xF0 != 0x20		// SCSI command
	return carry clear
if xdata 0x3018 != 0
	return carry clear
irrelevant; now we know data is already read by this point
```
So data is already read by this point, and thus we need to keep searching.

**ASIDE** - Hold up a sec. While writing FirmwareGeneralnotes.md, I came to a realization: The command bytes start at 0x300F. 0xF is the offset of the command bytes in a CBW packet. Is this firmware communicating using BBB after all, and just not checking for a signature?...

## Known boot ROM routines
I believe these are provided by the boot ROM; if they are actually in RAM and copied on system startup, I do not know (TODO).
- 0x1BBC - compare R0 to R4, R1 to R5, R2 to R6, R3 to R7
- 0x2B50 - memcpy
	- R6:R7 - (TODO source or destination? initial observation suggested source but after seeing random register assignment order I'm not sure (~0x821A suggests that it is indeed source))
	- R4:R5 - (TODO opposite of above)
	- R3 - length (bytes)
	- Is 0x2B6C the same, but for copying from code memory?

TODO could 0x329F be related to disk transfers — or worse, to encryption?

TODO could 0x2B46 be memclr?

TODO 0x1C2D does not return to caller; figure out how it processes the data after it (a word, the SCSI command arrays mentioned above; some more data)
