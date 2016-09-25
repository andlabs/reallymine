const NumSectorsAtATime = 102400

	// TODO make sure infile is not a device
	// we must outright forbid it because we aren't running sector-to-sector anymore
	// TODO or are we?

	// TODO decrypt a few sectors to verify the partition table

	fmt.Printf("Beginning decryption!\n")
	sectors := make([]byte, NumSectorsAtATime * SectorSize)
	n := int64(0)
	inmb := insize / 1024 / 1024
	for DecryptNext(in, out, bridge, c, sectors) {
		n += NumMBAtATime
		fmt.Printf("%d MB / %d MB complete.\n", n, inmb)
	}
