package poller

import (
	"io"
	"time"
)

type ScanPoller struct {
	dropInterval time.Duration
	batchMaxSize int
	scanningChan <-chan scanResult
}

func New(
	dropInterval time.Duration,
	batchMaxSize int,
	scanner inputScanner,
) *ScanPoller {
	return &ScanPoller{
		dropInterval: dropInterval,
		batchMaxSize: batchMaxSize,
		scanningChan: repeatScans(scanner),
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
				return []string{}, io.EOF
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
