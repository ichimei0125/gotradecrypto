package indicator

import "math"

// 计算平均值
func mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// 计算标准差
func standardDeviation(data []float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	avg := mean(data) // 计算平均值

	// 计算每个数据点与平均值的平方差
	sumOfSquares := 0.0
	for _, value := range data {
		diff := value - avg
		sumOfSquares += diff * diff
	}

	// 返回平方差的均值的平方根
	return math.Sqrt(sumOfSquares / float64(len(data)))
}
