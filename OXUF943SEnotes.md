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
