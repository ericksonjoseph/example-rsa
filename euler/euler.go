package euler

func Phi(i int64) int64 {
	return totient(i, getPrimes(1000000))
}

func getPrimeFactors(n int64, primes []int64) map[int64]int64 {
	primeFacts := make(map[int64]int64) // keeps track of prime factor : exponent pairs
	for n != 1 {
		for i := 0; i < len(primes); i++ {
			if n%primes[i] == 0 {
				val, ok := primeFacts[primes[i]]
				if !ok {
					val = 0
				}
				primeFacts[primes[i]] = val + 1
				n = n / primes[i]
				break
			}
		}
	}
	return primeFacts
}

func getPrimes(N int64) []int64 {
	isComposite := make([]bool, N)
	primes := []int64{}
	for i := int64(2); i < N; i++ {
		if !isComposite[i] {
			primes = append(primes, i)
			for x := i + i; x < N; x += i {
				isComposite[x] = true
			}
		}
	}
	return primes
}
func totient(n int64, primes []int64) int64 {
	primeFacts := getPrimeFactors(n, primes)

	ans := n

	for prime := range primeFacts {
		ans = ans * (prime - 1) / prime
	}
	return ans
}
