# reallymine: Western Digital MyBook/MyPassport decryption

`reallymine` is a program that decrypts the encrypted hard drives of Western Digital MyBook and MyPassport external hard drives (and some rebranded derivatives).

Currently, it can only decrypt JMicron and Initio bridge chip-based devices tht use AES-256-ECB encryption. I'd love to expand this to cover Symwave and PLX/Oxford Semiconductor bridge chips and the other known encryption modes, but I need your help; see below. It also does not currently handle entering passwords; if your drive is password-protected (and the bridge chip requires a password) but most of the work is already there (in `kek.go`); I just need to write the code that actually lets you type in a password, and then we'll be fine.

Simply run the program, providing the drive to decrypt and a file that the decrypted image will be stored to:

```
reallymine encrypted decrypted
```

`reallymine` **never overwrites a file that already exists**; by extension, it does not allow in-place decryption.

Note that I make no guarantees about whether running `reallymine` off an existhing hard drive will wear the drive out. It does not replace GNU ddrescue as a damaged-disk recovery tool. If in doubt, run GNU ddrescue first, then run `reallymine` off the rescued image.

## Contributing
As I mentioned earlier, `reallymine` is vastly incomplete. It only handles two of the four known bridge chips Western Digital used, and only supports one encryption mode. If you're willing to provide a few sectors from your drive (typically one of the last sectors and a few of the first ones), you can do so in the github issue tracker, and I can use them to improve this program! (Don't worry; I only need the boot sectors and decryption key; I won't need any of your actual data. The sectors won't go into the source repository either.)

## License
Because of those "data recovery experts" mentioned in notes/story.md, this project is licensed under the GPL version 3. You should be the one who owns your data, not other people. (In fact I'm wondering if this whole encryption thing is solely in place for their benefit.)

TODO should I switch to Affero GPL, just to be safe?

## Thanks (TODO)
- Xenesis (minor THUMB help)
- Sik (minor documentation fixes)
- FraGag (minor 68020 information)
- fd0 (irc.freenode.net #go-nuts; help with dealing with decryption keys)

## TODOs
- Elaborate on this README a bit; mention notes.
