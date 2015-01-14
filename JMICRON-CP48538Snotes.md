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
