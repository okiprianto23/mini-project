package transaction

import (
	"fmt"
	"time"
)

func (t transactionService) generateContractNumber(prefix string) (string, error) {
	// Dapatkan tanggal saat ini
	now := time.Now()
	datePart := now.Format("02012006") // Format ddmmyyyy

	// Increment angka urut di Redis
	key := "contract_number:" + prefix
	count, err := t.redis.Incr(key).Result()
	if err != nil {
		return "", err
	}

	// Format nomor kontrak
	contractNumber := fmt.Sprintf("%s-%04d", datePart, count)
	return contractNumber, nil
}
