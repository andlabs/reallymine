The Initio INIC-1607E is one of a family of Intel 8051-based USB-SATA bridge chips made by Initio. This seems to be the most common type of encryption chip in MyBooks, judging solely from the ratio of instances of people asking for encryption support with this vs. other manufacturers's chips. The E at the end signals encryption support; other entries in the family do not feature encryption. (Newer chips don't even use the 8051, instead having an unspecified "32-bit" CPU architecture.) Apart from promotional webpages, there doesn't seem to be any information about this chip (though I would need to check again; TODO).

I won't release the dump that I got from the same place Western Digital's own firmware updater software gets it from, nor will I release the disassembly — at least not on github. Here's what you need to know about the file, though:

```
filename: Apollo_R1-2018-20101026.bin
size:     65536 bytes
crc32:    24fc788b
md5:      0ac4728d7f7b9b75b936838432eef051
sha1:     f1fb04a4c652430dde51ab337d287287281d4660
```

As mentioned above, this chip is 8051-based.

A big problem is for whatever reason most of the lower half of the ROM is all 0xFF bytes! There's a standard `jmp @A+DPTR` jump table with the `jmp` instruction itself at ROM 0xABC2; it looks like this:

```
code:0000ABBC code_ABBC:                              ; CODE XREF: code:code_ABB7↑j
code:0000ABBC                 mov     DPTR, #0x2BC3
code:0000ABBF                 mov     R0, A
code:0000ABC0                 add     A, R0
code:0000ABC1                 add     A, R0
code:0000ABC2                 jmp     @A+DPTR
code:0000ABC3 ; ---------------------------------------------------------------------------
code:0000ABC3                 ljmp    code_2BE4
code:0000ABC6 ; ---------------------------------------------------------------------------
code:0000ABC6                 ljmp    code_2C24
code:0000ABC9 ; ---------------------------------------------------------------------------
code:0000ABC9                 ljmp    code_2C3F
code:0000ABCC ; ---------------------------------------------------------------------------
code:0000ABCC                 ljmp    code_2C4A
code:0000ABCF ; ---------------------------------------------------------------------------
code:0000ABCF                 ljmp    code_2C62
code:0000ABD2 ; ---------------------------------------------------------------------------
code:0000ABD2                 ljmp    code_2C7C
code:0000ABD5 ; ---------------------------------------------------------------------------
code:0000ABD5                 ljmp    code_2CAC
code:0000ABD8 ; ---------------------------------------------------------------------------
code:0000ABD8                 ljmp    code_2D1C
code:0000ABDB ; ---------------------------------------------------------------------------
code:0000ABDB                 ljmp    code_2D27
code:0000ABDE ; ---------------------------------------------------------------------------
code:0000ABDE                 ljmp    code_2D40
code:0000ABE1 ; ---------------------------------------------------------------------------
code:0000ABE1                 ljmp    code_2D51
code:0000ABE4 ; ---------------------------------------------------------------------------
```

Loading the ROM as-is results in each of these `ljmp`s jumping into the 0xFF block. The `mov DPTR` line, combined with these jump target addresses, implies that this code should be in the 0x2000 region instead. So let's see what happens if we load the lower half of the ROM at address 0x0 instead.

From this point on, addresses prefixed with a ~ should be with this mapping.

At ~0x2682 is a list of device ID mappings:
```
type xxxxTODO struct {
	vendorID		uint16
	productID		uint16
	unknown		uint16
	someString	*string	// ASCII; null-terminated (TODO match with JMicron notes)
	someString2	*stirng	// ASCII; null-terminated
	someString3	*string	// ASCII; null-terminated
}
```
(where pointers are two bytes wide); the last entry in this array begins at ~0x2706. Strangely, this last entry refers to a product called the "Western Digital Bit Bucket", vendor ID 0x1058 product ID 1.

The function at ~0x3114 writes R4, R5, R6, and R7 (in that order) to successive bytes at DPTR.

RAM 0x7D20 seems to be where the key sector is written to when creating...?

~0x72BE is memcpy(). R4:R5 is the destination, R6:R7 is the source, and R3 is the byte count.

The function to handle SCSI commands seems to be ~0x315D. The list of commands to parse is an array taken from the instruction stream of the form:
```go
; TODO fix formatting
type xxxTODO struct {
	code			pointer		// 16-bit
	command		byte
}
```
with the list ending with a zero `code`. After this zero `code` is the address to jump to if no command matches. In other words, a call that would only handle the READ(10) and READ(16) commands would be
```
	; TODO fix formatting
	mov	A, commandByte
	lcall	TODO
	struct xxxx <code_READ10, SCSI_READ_10>
	struct xxxx <code_READ16, SCSI_READ_16>
	.word 0
	.word code_neither_READ10_nor_READ16
```
where `SCSI_READ_10` and `SCSI_READ_16` are the command byte values themselves.

These SCSI command lists seem to be separated by overall task; the code that is run on a read operation appears to be at ~0x2331, with the list itself starting at ~0x2337.
