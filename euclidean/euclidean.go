package euclidean

import (
	"math/big"

	"github.com/sirupsen/logrus"
)

func Get(e, phi *big.Int) *big.Int {

	x := big.NewInt(phi.Int64())
	y := big.NewInt(phi.Int64())
	f := big.NewInt(1)
	var i int

	for e.Int64() != 1 {
		logrus.Tracef("Column 1 %d over %d\n", x, e)
		logrus.Tracef("divide %d by %d\n", x, e)

		newX := big.NewInt(e.Int64())
		var mul *big.Int
		mul, e = new(big.Int).DivMod(x, e, &big.Int{})
		logrus.Tracef("= %d r%d\n", mul, e)

		x = newX

		logrus.Tracef("Column 2 %d over %d\n", y, f)
		logrus.Tracef("multiply %d by %d\n", f, mul)
		p := new(big.Int).Mul(f, mul)
		logrus.Tracef("= %d\n", p)
		diff := new(big.Int).Sub(y, p)
		logrus.Tracef("sub %d - %d = %d\n", y, p, diff)
		if diff.Int64() < 0 {
			diff.Mod(diff, phi)
			logrus.Tracef("diff %d\n", diff)
		}
		if e.Int64() == 1 {
			return diff
		}
		y = big.NewInt(f.Int64())
		f = big.NewInt(diff.Int64())
		if i == 2 {
			break
		}
		i++
	}

	return nil
}
