package conv

import (
		"math/big"
)

func BigFloatFromStr(str string) *big.Float {
	var b big.Float
	b.SetString(str)
	return &b
}


