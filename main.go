package main

import (
	"fmt"
	"math"
	"math/big"

	"bitbucket.org/erickson1/rsa-example/euclidean"
	"bitbucket.org/erickson1/rsa-example/euler"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Printf("PHI(n): %d\n", euler.Phi(15))
	for _, n := range []int64{3} {
		if testEulersTheorem(2, n) {
			fmt.Printf("WORKS: %d\n", n)
		}
	}
	fmt.Printf("---------\n")
	tutorial()
	workingExample()
}

// testEulersTheorem m ^ ø(n) ≈ 1 mod n
// M to the power of Phi N is congruent
// to 1 mod n as long as M and N do not
// share a common factor (co-prime)
//
// Consider the numbers 2 and 3
//
//  ^1^2  ^3      ^4              ^5                              ^6
// | | |   |       |               |                               |
// |_._.___._______._______________|_______________________________|
// |__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|
// The counts will never intersect
// (they do not share a common factor)
func testEulersTheorem(m, n int64) bool {

	always1 := math.Mod(
		math.Pow(
			float64(m),
			float64(euler.Phi(n)),
		),
		float64(n),
	)
	if always1 != 1 {
		gcd := (&big.Int{}).GCD(nil, nil, big.NewInt(m), big.NewInt(n)).Int64()
		fmt.Printf("%d and %d share a common factor of %d\n", m, n, gcd)
		return false
	}

	return true
}

func tutorial() {

	// Calculating the PHI of any number is hard
	// but if you know the prime factors of that
	// number its easy.
	// e.g. phi of the following number
	n := 77
	fmt.Printf("Phi of %d = %d\n", n, euler.Phi(int64(n)))

	// But when it comes to prime numbers, we can calculate
	// it easily because the phi of any prime number (x) is that
	// number minus 1 i.e. PHI(x) = x - 1
	fmt.Printf("Phi of %d = %d\n", 7, euler.Phi(7))
	fmt.Printf("Phi of %d = %d\n", 11, euler.Phi(11))

	// So if we know the prime factors of n, we can
	// quickly get the phi of those numbers and multiply
	// them together to the the phi of n as long as the
	// factor numbers are prime numbers!
	fmt.Printf("phi(%d) = phi(%d) * phi(%d)\n", n, 7, 11)

	// If we don't know the prime factors of n before hand
	// we will need to calculate it, which can take years
	factors := primeFactors(int64(n))
	fmt.Printf("1000 years later... Prime Factorization of %d is %+v\n", n, factors)

	// Now link this to modular exponentiation (clock arithmetic)

	// Pick 2 numbers that do not share a common factor
	n = 8
	m := float64(5)
	exponent := float64(euler.Phi(int64(n)))

	result := math.Pow(m, exponent)
	fmt.Printf("m ^ phi(n) = %f\n", result)

	always1 := math.Mod(result, float64(n))
	fmt.Printf("m ^ phi(n) mod n = %f\n", always1)

	if always1 != 1 {
		// n & m cannot share any prime numbers
		fmt.Printf("Error: you picked bad numbers\n")
		return
	}

	// Breakthrough
	// m ^ k*phi(n)+1 ~ m mod n
	k := int64(2)
	fmt.Printf("\nFinal function:\nm ^ k*phi(n)+1 ~ m mod n\n____\n")
	a := math.Mod(math.Pow(m, float64(k*euler.Phi(int64(n))+1)), float64(n))
	b := math.Mod(m, float64(n))
	fmt.Printf("m ^ k*phi(n)+1 mod n = %f\n", a)
	fmt.Printf("m mod n = %f\n", b)
	if a != b {
		// You chose a bad k
		fmt.Printf("Bad values\n")
		return
	}

	fmt.Printf("Done showing examples\n\n")
}

func workingExample() {

	pk, sk := generateKeyPair()
	fmt.Println("Public Key", pk)
	encrypted := encrypt(big.NewInt(int64(89)), pk)
	fmt.Println("Encrypted with public key", encrypted)
	decrypted := decrypt(encrypted, sk)
	fmt.Println("Decrypted with private key", decrypted)
}

// generateKeyPair generates a public and private key pair
func generateKeyPair() (publicKey, privateKey) {

	logrus.SetLevel(logrus.DebugLevel)
	fmt.Printf("Generating key pair")

	// Choose two random prime numbers
	// of similar length @TODO why doesnt 5 & 11 work with e = 3 ?
	prime1 := int64(4817)
	prime2 := int64(3821)
	fmt.Printf("Starting with 2 prime numbers %d and %d\n", prime1, prime2)

	if !(big.NewInt(prime1)).ProbablyPrime(100) {
		panic(fmt.Sprintf("%d is not prime", prime1))
	}
	if !(big.NewInt(prime2)).ProbablyPrime(100) {
		panic(fmt.Sprintf("%d is not prime", prime2))
	}

	// Get the prime product. recommended 4096 bits
	primeProduct := prime1 * prime2
	fmt.Printf("Getting prime product i.e. %d * %d\n", prime1, prime2)
	fmt.Printf("Prime Product = %d\n", primeProduct)

	// It will be hard for anyone to
	// compute the primes given only the product
	//fmt.Printf("... primes = %+v\n", primeFactors(primeProduct))

	// We can easily calculate Phi of N since
	// we know the factorization of n
	phiN := (prime1 - 1) * (prime2 - 1)
	if phiN != euler.Phi(primeProduct) {
		panic("failed to start with real prime numbers")
	}

	fmt.Printf("Calculated phiN easily = %d\n", phiN)

	// Now we pick a random e with the conditions that
	// it is a number > 2 and does NOT share a factor
	// with phiN or the primeProduct
	exponent := int64(3)
	fmt.Printf("e = %d\n", exponent)
	if exponent < 3 {
		panic("e must be > 2")
	}
	for _, i := range primeFactors(exponent) {
		for _, j := range append(primeFactors(phiN), primeFactors(primeProduct)...) {
			if j == i {
				panic(fmt.Sprintf("%d shares a factor of %d with %d or %d", exponent, i, phiN, primeProduct))
			}
		}
	}

	//Calculate d by using Euclidean's extended division
	//with phiN and e

	d := euclidean.Get(big.NewInt(exponent), big.NewInt(phiN)).Int64()
	fmt.Printf("d = %d\n", d)
	// d is the number that will undo the effect of e
	// i.e. (m^e mod n)^d mod n = m .. or .. d*e mod phiN = 1

	// Now you can calculate k
	// d = (k*phiN + 1) / exponent
	//k := (d*exponent - 1) / phiN
	//fmt.Printf("k = %d\n", k)
	// Notice that you can change k and it will still work
	// as long as d remains a solid number

	// Now hide everything except n and e
	return publicKey{primeProduct, exponent}, privateKey{primeProduct, d}
}

func encrypt(msg *big.Int, pk publicKey) *big.Int {
	logrus.Debugf("encrypting with exponent %d and n %d", pk.exponent, pk.primeProduct)
	return (&big.Int{}).Exp(
		msg,
		big.NewInt(pk.exponent),
		big.NewInt(pk.primeProduct),
	)
}

func decrypt(msg *big.Int, sk privateKey) *big.Int {
	logrus.Debugf("decrypting with exponent %d and n %d", sk.d, sk.primeProduct)
	return (&big.Int{}).Exp(
		msg,
		big.NewInt(sk.d),
		big.NewInt(sk.primeProduct),
	)
}

type publicKey struct {
	// primeProduct is the product of two
	// prime numbers. It will be hard to
	// find the Phi() of this number
	// unless you know the prime numbers that
	// were used to create it
	primeProduct int64
	// e is the exponent
	exponent int64
}

type privateKey struct {
	primeProduct int64
	d            int64
}

// Get all prime factors of a given number n
func primeFactors(n int64) (pfs []int64) {
	// Get the number of 2s that divide n
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n = n / 2
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	for i := int64(3); i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			pfs = append(pfs, i)
			n = n / i
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n > 2 {
		pfs = append(pfs, n)
	}

	return
}
