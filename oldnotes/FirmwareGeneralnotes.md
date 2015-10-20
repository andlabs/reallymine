## General Firmware Notes

Previously, each MyBook had its own firmware update software which came with the firmware and code to perform the update. Sometime recently (as of 2010, possibly), there is now one single firmware update program, which downloads both the firmware and code to perform the updating from the Internet.

I have not dived too deeply into the Windows version of this new updater, as the Mac OS X version, because of how the Apple Objective-C runtime works, is effectively unstripped (as far as Objective-C classes and methods go, anyway). There are two XML files at play here. The first one, which is also provided with the updater, contains information about each supported MyBook family device. This is called the [TODO get name]. The second XML file, which is always downloaded (or so it seems...), contains the filenames of all the firmwares and updater software, for both Windows and Mac. This is (somewhat confusingly) called the "device table".

This latter file is encrypted using a very basic xor encryption with a constant, looping block of 513 bytes (yes, you read that correctly). (The version I looked at, [TODO get version], had the function that did the decrypting unstripped as well...)

The device table I have is listed as version 3.2.5.3 and is dated 26 August 2014. This file conveniently begins with a changelog and has a brief comment explaining the bridge board and chip before each device's entry (though not always in full detail).

The code that actually performs the firmware updating is called a "device plugin". I have looked at some of them, but have not gathered much information out of them.

Filenames in the device table are stored as relative paths; they are relative to a single HTTP server that also provides the device table itself. The URL of the server is stored in one of the `plist` files in the OS X version.

## Manufacturers and basic chip information

The following five companies have manufactured USB-SATA bridge chips for the MyBook series:

* **ASMedia**<br>Variety of [TODO get architecture] chips. [TODO get model]
* **Initio**<br>Variety of 8051-based chips. The firmware I'm looking at is identified as being from the INIC-1607E.
* **JMicron**<br>Variety of 8051-based chips. The firmware I'm looking at is identified as being from the JMICRON-CP48 538S (though both ASMedia and Initio are also listed as having chips called "CP48"...).
* **PLX** ([formerly](http://www.bloomberg.com/apps/news?pid=newsarchive&sid=aEeIQGrHLbrI) **Oxford Semiconductor**)<br>Variety of ARM-based chips. The firmware I'm looking at is identified as being from the OXUF943SE.
* **Symwave**<br>Variety of Motorola 68020-based chips. No specific model numbers identified.

I have chosen to look at one firmware from each of these companies. The selected firmware should be from a known device that uses encryption. In cases where a manufacturer has made chips for MyBooks without encryption (Initio and Symwave, for instance), the device table will usually indicate that device as being unencrypted. The `*notes.md` files in this directory of the repository contain my individual notes.

**As I mentioned in the README, I won't release files on github. I will release whatever I need to release (if I need to release) on my own server, or on a friend's server if that friend allows.**
