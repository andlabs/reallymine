The PLX OXUF943SE is one of many different USB-SATA bridge chips used by WD MyBooks with encryption. The only bit of information about this chip that isn't covered by a NDA is the promotional booklet, which tells us that the firmware is ARM-based and that the chip has both AES-128 and AES-256 encryption hardware.

I won't release the dump that I got from the same place Western Digital's own firmware updater software gets it from, nor will I release the disassembly — at least not on github. Here's what you need to know about the file, though:

```
filename: Release-1014-20101007.bin
size:     257368 bytes
crc32:    8ca46984
md5:      9fdfaedc612e6d031f2bf05aeb6aa44a
sha1:     81158da148ac04b4719c2e0e835edbbbb8bfbace
```

Due to the nonavailability of documentation, I'll be working from scratch.

First thing to note is that the firmware primarily runs in THUMB mode for whatever reasn. ([Good thing I have](https://github.com/andlabs/idapyscripts) [experience with ARM-THUMB!](https://github.com/andlabs/mmbnmapdump))

The memory map appears to be simple:

```
start      end         purpose
0x00000000 0x3FFFFFFF? code
0x40000000 0x7FFFFFFF? unknown (hardware registers?)
0x80000000 ??????????  RAM
	(TODO is the RAM size 0x4000 bytes? stack is near that location)
```

If you're going to go at this with IDA, here are some tips:

- There will be a LOT of calls to registers, which isn't possible in a single instruction in THUMB. Expect to see lots of
```
	LDR     R3, =(sub_target+1)
	BL      sub_nearby
	B       loc_nearby
sub_nearby
	BX      R3
loc_nearby
	...
```
with the usual IDA comment fluff.
- Sometimes, IDA won't disassemble a subroutine called this way (or through one of the various jump tables). You will need to select its byte and the byte prior (remember, this is THUMB), press `c`, and choose Force from the confirmation dialog to get it to work. Fortunately, these functions are simple, but at least one seems to be vital (having the USB product ID!).
- Don't read the YY in `[SP,#0xNN+var_YY]` as increasing; these are actually negative numbers. To get a better idea of how local variables are laid out, press `q` over the `var_YY` to turn it back into the actual offset from `SP` (which doubles as a frame pointer).

Known RAM addresses for jumping:

(TODO)

## Attempts at figuring out how/where USB input is processed

```
code for filling in device descriptor - first entry of list 0x3AF78
which is a jump table called near the beginning of the (very long) function 0x13794
turns out this is the function to handle GET_DESCRIPTOR
the jump table is for the various high bytes of the wValue
argument R0 to this function holds a pointer to the location of the input buffer ([R0] is the address of the input buffer itself)

this is called from 0x13CB8, which is responsible for handling standard USB requests
again, argument R0 holds a pointer to the location of the input buffer

this is called from 0x10C2C (specifically at 0x10CAC)
it turns out this function locally allocates room for the buffer, giving 0x24 bytes of its stack and starting the buffer at offset 0x1C
which appears to be copied from R0 (input) + 0x14...?

with R0 pointing directly to the input buffer, 0x15BC8 extracts the two request type bits (byte 0 bits 5 and 6)
likewise 0x15BD0 handles the recipient (byte 0 bits 0-4)

0xF7C0 seems to be the routine that handles class requests /to the interface/ /for specific requests only/
likewise, 0x12B64 seems to be the routine that handles vendor-specific requests

now for 0xF7C0 requests:
- 0xF86E - BBB Get Max LUN (0xFE)
- 0xF7EE - BBB Reset (0xFF)

I'll run under the assumption that the MyBooks are literally bulk-only devices and thus follow the [Bulk-Only Specification](http://www.usb.org/developers/docs/devclass_docs/usbmassbulk_10.pdf).

There is only one instance of the BBB CBW signature, at 0xFCF4.
It is loaded at 0xFB64, with the value to compare with coming immediately prior.
The code goes all over the place after checking structure integrity, but it eventually stores the read/write flag in several places (and in several ways?) and copies the command bytes in the loop starting at 0xFC3A.

Two CSW signatures:
first is at 0xF05C; loaded at 0xEF72
second is at 0x105FC; loaded at 0x10548
the former seems to only be used for failures
the latter /might/ be the general purpose case?

anyway...
0x10C2C is called at 0x1071C (through one of those `BL`/`BX` things mentioned above)
this is part of function 0x1060C... yet again using R0
THANKFULLY this is where our journey ends, as this function is called (at 0x38312) with R0 set to 0x800004B0
this setting is done by the function 0x10A40 (which only sets R0 and nothing else)
the offset 0x14 mentioned earlier appears to be a pointer to the buffer
0x80001380 + 0x19 appears to be a flag that indicates whether or not we should actually do the processing...
the flag is set at 0x24686 to the low byte of R0 passed to function 0x24648 (after some more complex logic)

0x46000000 might be the location of the USB communication ports
0x1B8B0 might be the function that actually performs reading
0x1B854 might be the function that actually performs writing
```
