package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func PrintHex(pre string, v *big.Int) {
	fmt.Printf("%v = 0x%v\n", pre, hex.EncodeToString(v.Bytes()))
}

func PrintDec(pre string, v *big.Int) {
	fmt.Printf("%v = %v\n", pre, v.String())
}
