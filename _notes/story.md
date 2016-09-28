## The Story
I have two 1TB Western Digital MyBooks.

* Drive A, from sometime between 2009 and 2011
* Drive B, from the summer of 2012

Drive A was simply a backup drive for my previous computer, an iMac.

Drive B, however, was bought both to back up my previous laptop's first hard drive when it started failing, my MobileMe files when they discontinued iDisk, and my main OS drive until I replaced the laptop's internal drive.

The story here begins with Drive B; over the next few months, the power adapters from both Drive B *and* Drive A became flimsier until neither could power Drive B. As this was in the fall of 2012, during school, I was in a panic and tried whatever adapter I could find, thinking the big box around the actual plug would regulate voltage. Haha, yeah right: the first one I found that fit overloaded the drive. When I was finally able to perform a ddrescue on the drive... it came up as garbage data. Mostly, anyway; there were some Western Digital fiiles near the end of the drive. Fearing I had fried the drive but with the WD files serving as a sort of hope spot, I shelved the drive for a while.

Then, in January 2015, a friend needed a file that I knew predated the data lost to Drive B. So I took Drive A out of its case and plugged it into a USB chasis... and nothing happened. I did a hexdump of the drive itself and that came up as garbage too! Thinking my chasis was damaged, I plugged the iMac's internal HDD in... and it worked.

It didn't take much Googling to confirm what I subsequently suspected.

## The Facts
Several families of Western Digital MyBooks (and the portable equivalents, MyPassports), as well as several rebranded versions of such (some by HP, for instance), have mandatory, transparent, full-disk encryption. The encryption is performed by a chip on the USB-SATA bridge board.

The encryption is known to be standard AES-128 and/or AES-256, with no additional block ciphers (ECB) (though CBC and XCB models might exist as well). The chips vary between drives. In some cases, swapping disks between cases/bridge boards of the same model does work to decrypt the data, but this is not always the case.

The encryption chip also chops off the upper portion of the disk (or so). This portion, whose size I am not sure about, is the source of the WD files I mentioned earlier: it is actually a CD image that the bridge chip firmware exposes to the host OS (Windows or OS X) as a regular CD. The CD normally contains the program which unlocks the HDD if you gave it a password using the WD SmartWare utility (which is NOT the same thing as a regular ATA password). I do not think the password has any bearing on encryption (and the drive is still encrypted even without a password).

As it turns out, the encryption key isn't necessarily stored on the bridge chip. Instead, it's stored in two places: a "module" of the drive's "system area" (I don't know what this means, nor can I yet find a Linux utility that examines this - TODO) and as a backup in a sector near the end of the drive. This "key sector" contains several bits of information (notably the size of the drive that the bridge chip exposes).

## Cracking the Code
I spent much of the first few months of 2015 on independent research, then took a hiatus to focus on [other projects](https://github.com/andlabs/libui). You can see the results of this early research in the folder notes/old/. My research was done entirely by reverse-engineering firmware and Western Digital's VCD software. The firmware was downloaded from Western Digital's servers, based on reverse-engineered firmware upgraders.

In the meantime, three security researchers, Gunnar Alendal, Christian Kison, and modg, independently performed their own research, using hardware tools as well as software tools. Their paper, ["got HW crypto?: On the (in)security of a Self-Encrypting Drive series"](http://eprint.iacr.org/2015/1002.pdf), was published in September 2015, but I only found out a month later via Twitter. Their work went above and beyond what I ever did, to the point that **almost everything we need to recover a drive is finally public knowledge**. And now, with a little bit more figuring out, I can finally write the actual reallymine program. My research is thus now abandoned, as it is no longer needed; it is still available (as mentioned above).
