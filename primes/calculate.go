package primes

// this is an implementation of the Sieve of Eratosthenes!
func Calculate(max int) {
	marked := make(map[int]bool)

	p := 2 // let p = the first prime number

	// starting from p-squared, count up in increments of p and
	// mark each of these numbers which is greater than p-squared.

	// repeat this process for all numbers greater than p which have not been marked

	for j := p; j <= max; j++ {
		if _, ok := marked[j]; ok {
			continue // this number is already marked
		}

		for i := j; i <= max; i += j {
			if i >= (j * j) {
				marked[i] = true
			}
		}
	}

	for i := p; i <= max; i++ {
		if _, ok := marked[i]; ok {
			continue
		}
	}
}
