# reallymine: Western Digital MyBook/MyPassport decryption

`reallymine` is a program that decrypts the encrypted hard drives of Western Digital MyBook and MyPassport external hard drives (and some rebranded derivatives).

Currently, it can only decrypt JMicron, Initio, and Symwave bridge chip-based devices that use AES-256-ECB encryption. I'd love to expand this to cover PLX/Oxford Semiconductor bridge chips and the other known encryption modes, but I need your help; see below.

The program is command-based, with two main commands (one to get the decryption key and one to decrypt the drive automatically) and several helper commands that will facilitate research in expanding reallymine. The general usage is thus

```
$ reallymine [options] command [args...]
```

Pass `--help` for more detailed explanations.

## Installing
Stable versions of reallymine are available from the Releases page on GitHub.

reallymine is written in Go. If you want to build it from source, install Go and then simply run

```
$ go get github.com/andlabs/reallymine
```

This will get reallymine and its dependencies and place the resultant binary in your `$GOPATH/bin`.

If you want to manually download reallymine, you will need to have the dependencies installed separately:

```
github.com/mendsley/gojwe
	for the AES key-unwrapping code used to extract the DEK from Symwave chips
github.com/hashicorp/vault/helper/password
	for password entry
```

## Decrypting a Drive
The most common operation is decrypting an entire drive. Let's say the drive is at `/dev/sdb` and you want to decrypt it to a file `decrypted.img`. You would just say

```
$ reallymine decrypt /dev/sdb decrypted.img
```

reallymine will automatically find the sector on the drive that holds the encryption information, referred to by the program as the "key sector", and attempts to extract the decryption key without a password. If that fails, reallymine will ask you for a password. Once the right password is entered and the decryption key is extracted, reallymine will start decrypting the drive. This will take a while; sit tight.

`reallymine` **never overwrites a file that already exists**; by extension, it does not allow in-place decryption.

Note that I make no guarantees about whether running `reallymine` off an existhing hard drive will wear the drive out. It does not replace GNU ddrescue as a damaged-disk recovery tool. If in doubt, run GNU ddrescue first, then run `reallymine` with the rescued image.

## Getting the Decryption Key
You may want to perform decryption yourself, or do other things with the decryption key. For that, use the `getdek` command.

```
$ reallymine getdek /dev/sdb
```

In addition to printing the type of bridge chip your drive uses and the encryption key, which reallymine calls the DEK, reallymine will also print the steps needed to properly decrypt the data on the drive for every AES cipher block. For example, an Initio drive will say

```
bridge type Initio
DEK: A5DC4A231E88162A7066B063C2C31F1BDF248AF53F4F86F432C9E5414F88D280
decryption steps: swaplongs decrypt swaplongs
```

indicating that you first need to reverse all the 4-byte groups in each 16-byte block, then decrypt with the DEK, and then reverse the 4-byte groups again to get the final data.

## Researching with reallymine
reallymine has several research-oriented commands built in in addition to the two above. When contributing, I may ask you to run these commands to find out more about your specific scenario. You may also run them yourselves.

First are `dumplast` and `dumpkeysector`. `dumpkeysector` will try to find and dump the key sector on your drive as it is stored on disk. If that fails to detect the key sector, you can try `dumplast`, which gets the last sector on the drive that isn't all zero bytes; we can look at it to see what key sector you have. Both have the same syntax

```
$ reallymine dumplast /dev/sdb outfile.bin
$ reallymine dumpkeysector /dev/sdb outfile.bin
```

Alternatively, you can use `-` as an output filename to perform a hexdump on standard output.

The `decryptkeysector` command is like the `dumpkeysector` command, except it also decrypts the key sector with the encryption key that is used to encrypt that specific sector, which reallymine calls the KEK. The KEK changes when you change your password; the DEK never changes. Consequently, the KEK is used to encrypt the DEK to ensure the DEK doesn't leak out.

`decryptkeysector` has the same form as `dumpkeysector`, except it takes a third argument to specify the KEK. This can be a hexadecimal string to use a specific KEK, or one of the following special values:

```
-real    - behave like the decrypt command
-askonce - ask for a password once and only use the resultant KEK
-onlyask - only ask for a password until the right one is used
-default - use the default KEK (no password) only
```

The DEK can likely be read out of the decrypted key sector.

The `dumpfirst` command, which takes the same form as the `dumplast` command, dumps the first few sectors of your hard drive without decrypting them. This will likely contain the partition map of your drive, allowing it to be used to verify that a DEK is correct without leaking any of your sensitive data.

But simply knowing the DEK is not enough; you also need to know how to transform the data before and after decrypting to get the data back out properly. This is done with the `decryptfile` command, which does not deal with a disk at all. It takes four parameters: an input file to decrypt (or `-` for standard input), an output file to decrypt to (or `-` for a hexdump to stdout), the DEK as a hexadecimal string, and then a space-delimited string containing the decryption steps, such as those shown in the example output of the `getdek` command. Use `--help` for a full list of possible steps.

More specific usage information can be seen with `--help`.

## Contributing
reallymine is already quite capable, but is still in need of improvement to handle every possible case. If your drive isn't handled already, feel free to open an issue on GitHub to contribute your key sectors and partition maps, either by following the steps above or with our help. (Don't worry; I only need the boot sectors and decryption key; I won't need any of your actual data. The sectors won't go into the source repository either.)

Code contributions are also welcome.

## License
This project is licensed under the GPL version 3. This is to ensure that the research that went into reallymine stays open.

TODO should I switch to Affero GPL, just to be safe?

## Thanks (TODO)
- Xenesis (minor THUMB help)
- Sik (minor documentation fixes)
- FraGag (minor 68020 information)
- fd0 (irc.freenode.net #go-nuts; help with dealing with decryption keys)
- Everyone else from IRC and the GitHub issues I forgot to thank
