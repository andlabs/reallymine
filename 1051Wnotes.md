The ASMedia 1051W (TODO some other chip name too?) is a USB-SATA bridge chip used by WD MyBooks. I can't easily tell if it has encryption or not, as it was possibly used on both ecnrypted and unencrypted MyBooks (with the unencrypted ones explicitly pointed out in the Western Digital firmware update tools' database); we'll just hope for the best.

I won't release the dump that I got from the same place Western Digital's own firmware updater software gets it from, nor will I release the disassembly — at least not on github. Here's what you need to know about the file, though:

```
filename: Release-CFU-1065-20140325.bin
size:     61440 bytes
crc32:    2607c88c
md5:      ce3f7458fa5d61269e553dccaa821fd5
sha1:     92d9a4a0b381d1fd0324587e317b4f857147e569
```

Due to the nonavailability of documentation (TODO re-verify), I'll be working from scratch.
