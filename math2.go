package goincv

import "math"

// OneEuroFilter is a filter that smooths a signal.
type OneEuroFilter struct {
	// The parameters.
	minCutoff float64
	beta      float64
	dCutoff   float64
	// Previous values.
	xPrev  float64
	dxPrev float64
	tPrev  float64
}

// NewOneEuroFilter creates a new OneEuroFilter.
func NewOneEuroFilter(t0, x0, dx0, minCutoff, beta, dCutoff float64) *OneEuroFilter {
	return &OneEuroFilter{
		minCutoff: minCutoff,
		beta:      beta,
		dCutoff:   dCutoff,
		xPrev:     x0,
		dxPrev:    dx0,
		tPrev:     t0,
	}
}

// smoothingFactor computes the smoothing factor.
func smoothingFactor(tE, cutoff float64) float64 {
	r := 2 * math.Pi * cutoff * tE
	return r / (r + 1)
}

// exponentialSmoothing computes the exponential smoothing.
func exponentialSmoothing(a, x, xPrev float64) float64 {
	return a*x + (1-a)*xPrev
}

// Filter filters the signal.
func (f *OneEuroFilter) Filter(t, x float64) float64 {
	tE := t - f.tPrev

	// The filtered derivative of the signal.
	aD := smoothingFactor(tE, f.dCutoff)
	dx := (x - f.xPrev) / tE
	dxHat := exponentialSmoothing(aD, dx, f.dxPrev)

	// The filtered signal.
	cutoff := f.minCutoff + f.beta*math.Abs(dxHat)
	a := smoothingFactor(tE, cutoff)
	xHat := exponentialSmoothing(a, x, f.xPrev)

	// Memorize the previous values.
	f.xPrev = xHat
	f.dxPrev = dxHat
	f.tPrev = t

	return xHat
}

// estimateAffine2D 使用最小二乘法估计 2D 仿射变换矩阵
func EstimateAffine2D(src, dst []float32) []float32 {
	// 将坐标系移到图像的中心
	srcMean := make([]float32, 2)
	dstMean := make([]float32, 2)
	for i := 0; i < len(src); i += 2 {
		srcMean[0] += src[i]
		srcMean[1] += src[i+1]
		dstMean[0] += dst[i]
		dstMean[1] += dst[i+1]
	}
	srcMean[0] /= float32(len(src)) / 2
	srcMean[1] /= float32(len(src)) / 2
	dstMean[0] /= float32(len(dst)) / 2
	dstMean[1] /= float32(len(dst)) / 2
	for i := 0; i < len(src); i += 2 {
		src[i] -= srcMean[0]
		src[i+1] -= srcMean[1]
		dst[i] -= dstMean[0]
		dst[i+1] -= dstMean[1]
	}

	// 计算协方差矩阵
	cov := make([]float32, 4)
	for i := 0; i < len(src); i += 2 {
		cov[0] += src[i] * dst[i]
		cov[1] += src[i] * dst[i+1]
		cov[2] += src[i+1] * dst[i]
		cov[3] += src[i+1] * dst[i+1]
	}

	// 计算特征值和特征向量
	a := cov[0] + cov[3]
	b := cov[2] - cov[1]
	c := float32(math.Sqrt(float64(a*a + b*b)))
	sin := b / c
	cos := a / c

	// 计算仿射变换矩阵
	return []float32{
		cos, -sin, srcMean[0] - cos*dstMean[0] + sin*dstMean[1],
		sin, cos, srcMean[1] - sin*dstMean[0] + cos*dstMean[1],
	}
}
