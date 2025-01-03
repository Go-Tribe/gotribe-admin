package util

// FenToYuan converts price from fen (int) to yuan (float64)
func FenToYuan(fen int) float64 {
	return float64(fen) / 100.0
}

// YuanToFen converts price from yuan (float64) to fen (int)
func YuanToFen(yuan float64) int {
	return int(yuan * 100)
}
