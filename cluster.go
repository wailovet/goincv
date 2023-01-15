package goincv

import (
	"github.com/mpraski/clusters"
)

type KMeansCluster struct {
	data [][]float64
	c    clusters.HardClusterer
}

func (k *KMeansCluster) Add(data []float32) {
	k.data = append(k.data, F32ToF64(data))
}

func (k *KMeansCluster) Learn(clasNum, iterations int) (err error) {
	k.c, err = clusters.KMeans(iterations, clasNum, clusters.EuclideanDistance)
	if err != nil {
		return err
	}
	err = k.c.Learn(k.data)
	return err
}

func (k *KMeansCluster) Predict(data []float32) int {
	return k.c.Predict(F32ToF64(data))
}
