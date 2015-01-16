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

The encryption is known to be standard AES-128, with no additional block ciphers (ECB). The chips vary between drives. In some cases, swapping disks between cases/bridge boards of the same model does work to decrypt the data, but this is not always the case.

The encryption chip also chops off the upper portion of the disk (or so). This portion, whose size I am not sure about, is the source of the WD files I mentioned earlier: it is actually a CD image that the bridge chip firmware exposes to the host OS (Windows or OS X) as a regular CD. The CD normally contains the program which unlocks the HDD if you gave it a password using the WD SmartWare utility (which is NOT the same thing as a regular ATA password). I do not think the password has any bearing on encryption (and the drive is still encrypted even without a password).

As it turns out, the encryption key isn't necessarily stored on the bridge chip. Instead, it's stored in two places: a "module" of the drive's "system area" (I don't know what this means, nor can I yet find a Linux utility that examines this - TODO) and as a backup in a sector near the end of the drive. This "key sector" contains several bits of information (notably the size of the drive that the bridge chip exposes).

## So what have I done so far?
So far I've found the key sector and attempted to brute force the key out of it. One thing's for sure: it is NOT stored directly as such in the sector.

So how is it stored?

I don't know. The only people who do know seem to be so-called "data recovery experts", who have [chosen not to reveal this information](http://forum.hddguru.com/viewtopic.php?t=21584) [lest it hurt their business](http://forum.hddguru.com/viewtopic.php?t=24567&f=1&start=0#p165906). **Bullshit**. I'm not trusting my data to strangers.

As a result, this project is licensed under the GPL version 3. You should be the one who owns your data, not other people. (In fact I'm wondering if this whole encryption thing is solely in place for their benefit.) It's also appalling that there is only one website on the entire Internet dedicated to cracking this nut.

[This is not an unsolved problem; there are commercial utilities that can do what I am aiming to do here.](http://www.acelaboratory.com/news/newsitem.php?itemid=115) [Someone else has posted an encryption key from another peron's drive](http://forum.hddguru.com/viewtopic.php?t=19408&f=1&start=40&#p136073) [given only these screenshots of the key sector](http://forum.hddguru.com/viewtopic.php?t=19408&f=1&start=0#p131488). (A post on the third link indicates that the bridge chips appear to be Intel 8051 derivatives.)

I have attempted to reverse-engineer Unlock.exe. So far, I figured out that the SmartWare password is salted and hashed with a buggy implementation of SHA-256. I do not know what the salt is, but it is related to the drive being passworded. I do not think this will help figure out the encryption key.

The OS X equivalent of Unlock.exe is a bit more open: the salt and iteration count seem to be more variable and the default salt appears to be `WDC.`. No idea if this is a guarantee. No idea if the Windows version writes this as a UTF-8 or as a UTF-16 string. Still don't think this is related, even though the OS X version talks about encryption a lot.

Of note: the OS X equivalent of Unlock.exe, which is not stripped by virtue of the design of the Objective-C runtime, calls the block with the password hint the "handy store security block" and begins with the byte sequence `00 01 44 57` (the last two bytes being `WD` in reverse).

I have a dump of the "UF924DS" bridge chip firmware version r1.08a from 2007; this appears to be before WD started encrypting the drives, as there don't seem to be AES constants in the firmware (though I might not be looking hard enough). I am currently up to trying to find newer firmware versions.

## Contributions
**PLEASE**. Just be willing to share this knowledge for the sake of everyone who is unfortunate enough to own one of these drives.

## TODO
- brute force program: http://forum.hddguru.com/viewtopic.php?t=19408&f=1&start=40&#p136073 describes some necessary ciphertext alterations (byteswapping)
- get a more comprehensive list of the affected drives
	- confirm if MyPassport is affected /in exactly the same way/ and update the description of this repo
- my Drive A has two sectors that begin with `WD`, the key sector and some other, and neither of them use `WDv1` (and the part between the two `WD`s is 16 bytes smaller)
- drives that use some(?) Symwase bridge chips have an entirely different format for this sector; TODO
- get info about Drive A's bridge chip; I still have it and it still works
	- find Drive B's bridge chip; no idea if it can be salvaged but I still have it somewhere around here
- get the model number from Drive A's case
	- Drive B's case might be lost, however...
