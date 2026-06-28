package netclient

import (
	"bufio"
	"io"
	"strings"
)

// ScanKeyword - потоковый поиск keyword в теле HTTP-ответа
func ScanKeyword(body io.Reader, keyword string, limitBytes int64) (bool, error) {
	if keyword == "" {
		return false, nil
	}
	limited := io.LimitReader(body, limitBytes)
	scanner := bufio.NewScanner(limited)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, keyword) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
