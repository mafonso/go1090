package main

import (
	"fmt"
)

// https://adsb-decode-guide.readthedocs.io/en/latest/
func parseModeS(message []byte) {

	fmt.Printf("*%x;\n", message)

	msgType := uint((message[0] & 0xF8) >> 3)

	if msgType == 0 {
		fmt.Printf("DF 0: Short Air-Air Surveillance.\n")

		if uint((message[0] & 0x04)) == 1 {
			fmt.Printf("  VS             : Ground\n")
		} else {
			fmt.Printf("  VS             : Airborne\n")
		}

		fmt.Printf("  CC             : %d\n", uint((message[0]&0x02)>>1))
		fmt.Printf("  SL             : %d\n", uint((message[1]&0xE0)>>5))

	} else if msgType == 4 {
		fmt.Printf("DF 4: Surveillance, Altitude Reply.\n")
	} else if msgType == 5 {
		fmt.Printf("DF 5: Surveillance, Identity Reply.\n")
	} else if msgType == 11 {
		fmt.Printf("DF 11: All Call Reply.\n")
	} else if msgType == 16 {
		fmt.Printf("DF 16: Long Air to Air ACAS.\n")
	} else if msgType == 17 {
		fmt.Printf("DF 17: ADS-B message.\n")
	} else if msgType == 18 {
		fmt.Printf("DF 17: Extended Squitter.\n")
	} else if msgType == 19 {
		fmt.Printf("DF 19: Military Extended Squitter.\n")
	} else if msgType == 20 {
		fmt.Printf("DF 20: Comm-B, Altitude Reply.\n")
	} else if msgType == 21 {
		fmt.Printf("DF 21: Comm-B, Identity Reply.\n")
	} else if msgType == 22 {
		fmt.Printf("DF 22: Military Use.\n")
	} else if msgType == 24 {
		fmt.Printf("DF 24: Comm D Extended Length Message.\n")
	} else if msgType == 32 {
		fmt.Printf("SSR : Mode A/C Reply.\n")
	} else {
		fmt.Printf(" Unknown DF Format.\n")
	}

	fmt.Printf("\n\n")
}
