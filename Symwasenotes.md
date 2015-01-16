Western Digital's firmware update tool's database is nondescript about which Symwase chips the firmwares it does have map to (only saying "Symwase R2"), but does indicate whether a given firmware is used on an encrypted or an unencrypted device. Let's take an encrypted one, of course :V

I won't release the dump that I got from the same place Western Digital's own firmware updater software gets it from, nor will I release the disassembly — at least not on github. Here's what you need to know about the file, though:

```
filename: Release-1016_1U_Artemis-20110603.Bin
size:     89222 bytes
crc32:    794de485
md5:      85dfcaf27a934171f8838573ec4e01e6
sha1:     5b7f3f80795b4689a4c31f358f8fcdf00b2d7567
```

Examining hexdumps of the firmwares almost immediately pointed me toward something unusual: `4E75`... suspiciously familiar. As it turns out, these chips indeed are powered by the Motorola 68000, the CPU architecture I'm most familiar with.

(TODO find the game that had a TOO EASY sound clip in a really gruff low-pitch voice and put a clip of that here)

In reality it's of the 68020-and-newer variety with more features. Also it wasn't really that easy until I determined that the OXUF943SE used the USB Mass Storage Device Bulk-Only mode and found that this firmware does too.

This is NOT a 68000 boot ROM; it appears to be loaded at $4000000. Addresses beginning with ~ are in the offset form. Nevertheless, it is mostly code, and virtually every function in the code has a frame pointer; that is,
```
	link	a6,#nn
	...
	unlk	a6
	ret
```
(`4E75` is the `ret`.)

The code to handle a CBW is at the function at ~$400CDCA; it checks one character at a time for some reason. The firmware seems to /expect/
- up to 16 CBW command bytes, copied elsewhere for processing unconditionally after parsing the CBW structure
- at least one CBW command byte, checking its values and returning errors (I think) before even doing any additional processing (outside of parsing the CBW structure)

BIG OL' TODO: the code that copies the 16 CBW command bytes uses longword reads from an odd address; this is illegal on the vanilla 68000 (misaligned memory access). Does the 68020 allow it?
