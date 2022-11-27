package poller

import (
	"time"
)

type ScanPoller struct {
	dropInterval time.Duration
	batchMaxSize int
	scanningChan <-chan scanResult
	scanner      inputScanner
}

func New(
	dropInterval time.Duration,
	batchMaxSize int,
	scanner inputScanner,
) *ScanPoller {
	scanningChan := make(chan scanResult)
	close(scanningChan)

	return &ScanPoller{
		dropInterval: dropInterval,
		batchMaxSize: batchMaxSize,
		scanner:      scanner,
		scanningChan: scanningChan,
	}
}

func (b *ScanPoller) Poll() ([]string, error) {
	dropBatchTicker := time.NewTicker(b.dropInterval)
	defer dropBatchTicker.Stop()

	batch := make([]string, 0, b.batchMaxSize)

	for {
		select {
		case aScanResult, ok := <-b.scanningChan:
			if !ok {
				b.scanningChan = repeatScans(b.scanner)

				continue
			}

			if aScanResult.err != nil {
				return []string{}, aScanResult.err
			}

			batch = append(batch, aScanResult.text)

			if len(batch) == cap(batch) {
				return batch, nil
			}
		case <-dropBatchTicker.C:
			if len(batch) > 0 {
				return batch, nil
			}
		}
	}
}
