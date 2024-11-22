package calculator

type VolumePriceCalculator struct {
	volumes []uint64
	weights []float64
}

func NewPriceCalculator(volumes []uint64, weights []float64) *VolumePriceCalculator {
	return &VolumePriceCalculator{
		volumes: volumes,
		weights: weights,
	}
}

func (v *VolumePriceCalculator) CalculateWeightedMovingAverage() float64 {
	var volumes = v.volumes
	var weights = v.weights
	if len(volumes) < len(weights) {
		weights = weights[:len(volumes)]
	}
	var totalVolume, totalWeight float64
	for i := len(volumes) - len(weights); i < len(volumes); i++ {
		weight := weights[i-(len(volumes)-len(weights))]
		totalVolume += float64(volumes[i]) * weight
		totalWeight += weight
	}
	return totalVolume / totalWeight
}

func (v *VolumePriceCalculator) CalculateSimpleMovingAverage() float64 {
	var volumes = v.volumes
	var n = len(volumes)
	if len(volumes) < n {
		n = len(volumes)
	}
	var totalVolume uint64
	for i := len(volumes) - n; i < len(volumes); i++ {
		totalVolume += volumes[i]
	}
	return float64(totalVolume) / float64(n)
}
