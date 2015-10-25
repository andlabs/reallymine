// 22 october 2015
package main

import (
	"fmt"
	"os"
)

func main() {
	RealMain()
//	QuickTestMain()
}

func errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func die(format string, args ...interface{}) {
	errf(format, args...)
	errf("\n")
	os.Exit(1)
}

func usage() {
	errf("usage: %s encrypted decrypted\n", os.Args[0])
	errf("	encrypted must exist; should not be a device\n")
	errf("	decrypted must NOT exist\n")
	os.Exit(1)
}

func RealMain() {
	if len(os.Args) != 3 {
		usage()
	}
	infname := os.Args[1]
	outfname := os.Args[2]

	in, err := os.Open(infname)
	if err != nil {
		die("error opening encrypted file %s: %v", infname, err)
	}
	defer in.Close()

	// TODO make sure infile is not a device

	insize, err := in.Seek(0, 2)
	if err != nil {
		errf("error finding size of encrypted file %s: %v", infname, err)
	}

	fmt.Printf("Finding key sector...\n")
	keySector, bridge := FindKeySectorAndBridge(in, insize)
	if bridge == nil {
		errf("Sorry, we couldn't find the key sector.\n")
		errf("Either the drive isn't a complete image,\n")
		errf("or the encryption isn't supported yet.\n")
		os.Exit(1)
	}
	fmt.Printf("Found %s.\n", bridge.Name())
	if !bridge.NeedsKEK() {
		fmt.Printf("You will not need to enter your password\n")
		fmt.Printf("for this bridge chip.\n")
	} else {
		fmt.Printf("Trying without a password...\n")
	}

	c := TryGetDecrypter(keySector, bridge, func(firstTime bool) (password string, cancelled bool) {
		if firstTime {
			fmt.Printf("The drive's password is needed to decrypt your drive.\n")
			fmt.Printf("Please enter it now.\n")
		} else {
			fmt.Printf("Password incorrect.\n")
		}
		// TODO
		os.Exit(2)
		panic("unreachable")
	})
	if c == nil {
		fmt.Printf("User aborted operation.\n")
		os.Exit(1)
	}

	// TODO decrypt a few sectors to verify the partition table

	_, err = in.Seek(0, 0)
	if err != nil {
		die("error seeking back to start of decrypted file %s: %v", infname, err)
	}

	out, err := os.OpenFile(outfname, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			errf("Error creating decrypted file %s: %v\n", outfname, err)
			errf("%s will not overwrite a file that already exists.\n", os.Args[0])
			errf("In particular, %s does not allow in-place decryption.\n", os.Args[0])
			os.Exit(1)
		}
		die("error creating decrypted file %s: %v", outfname, err)
	}

	fmt.Printf("Beginning decryption!\n")
	insize01p := (insize / SectorSize) / 1000
	n := int64(0)
	p := 0.0
	for {
		more := DecryptNextSector(in, out, bridge, c)
		if !more {
			break
		}
		n++
		if n == insize01p {
			n = 0
			p += 0.1
			fmt.Printf("%.1f%% complete.\n", p)
		}
	}

	fmt.Printf("Completed successfully!\n")
}

func QuickTestMain() {
	f, _ := os.Open(os.Args[1])
	fout, _ := os.Create(os.Args[2])

	size, _ := f.Seek(0, 2)
	keySector, bridge := FindKeySectorAndBridge(f, size)
	if keySector == nil {
		fmt.Println("no key sector found")
		return
	}
	fmt.Println("found " + bridge.Name())

	c := TryGetDecrypter(keySector, bridge, func(firstTime bool) (password string, cancelled bool) {
		if firstTime {
			fmt.Println("We need the drive's password to decrypt your drive.")
		} else {
			fmt.Println("Password incorrect.")
		}
		// TODO
		return "abc123", false
	})
	if c == nil {
		fmt.Println("User aborted.")
		return
	}

	_, err := f.Seek(0, 0)
	if err != nil {
		// TODO
		panic(err)
	}
	for DecryptNextSector(f, fout, bridge, c) {
		// TODO
	}
}
