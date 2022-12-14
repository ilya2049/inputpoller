package poller

type inputScanner interface {
	ScanInput() (string, error)
}

type scanResult struct {
	text string
	err  error
}

func repeatScans(scanner inputScanner) <-chan scanResult {
	c := make(chan scanResult)

	go func() {
		for {
			input, err := scanner.ScanInput()

			c <- scanResult{
				text: input,
				err:  err,
			}

			if err != nil {
				close(c)

				break
			}
		}
	}()

	return c
}
