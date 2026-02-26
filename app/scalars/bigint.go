package scalars

import (
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func UnmarshalBigNumber(v interface{}) (*big.Int, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("BigNumber must be a numerical string")
	}
	bigInt := new(big.Int)
	_, ok = bigInt.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid BigNumber value: %q", s)
	}

	return bigInt, nil
}

// MarshalBigNumber converts a *big.Int into a scalar value of type BigNumber.
func MarshalBigNumber(number *big.Int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		quotedString := strconv.Quote(number.String())
		_, _ = w.Write([]byte(quotedString))
	})
}
