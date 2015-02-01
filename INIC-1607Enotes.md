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
