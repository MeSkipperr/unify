package utils

func BytesPerSecToMbps(bytesPerSec int64) float64 {
	return float64(bytesPerSec*8) / 1_000_000
}
