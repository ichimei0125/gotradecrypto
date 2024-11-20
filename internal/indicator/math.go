package indicator

import (
	"math"
	"sort"
)

func calculateSMA(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

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

func maxFloat64(slice []float64) float64 {
	sort.Sort(sort.Float64Slice(slice))
	return slice[len(slice)-1]
}

func minFloat64(slice []float64) float64 {
	sort.Sort(sort.Float64Slice(slice))
	return slice[0]
}
