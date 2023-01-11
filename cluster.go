package goincv

import (
	"encoding/json"
	"math"
)

type KMeansCluster struct {
	data      [][]float32
	centroids [][]float32
}

func (k *KMeansCluster) Add(data []float32) {
	k.data = append(k.data, data)
}

func (k *KMeansCluster) Learn(clusters, iterations int) {
	// Initialize centroids randomly
	for i := 0; i < clusters; i++ {
		centroid := k.data[i]
		k.centroids = append(k.centroids, centroid)
	}

	for i := 0; i < iterations; i++ {
		// Assign points to nearest centroid
		clusters := make([][][]float32, clusters)

		// find nearest centroid for each point
		for _, point := range k.data {
			nearestCentroid := 0
			minDistance := float32(math.MaxFloat32)
			for j, centroid := range k.centroids {
				distance := euclideanDistance(point, centroid)
				if distance < minDistance {
					minDistance = distance
					nearestCentroid = j
				}
			}
			clusters[nearestCentroid] = append(clusters[nearestCentroid], point)
		}

		// Recalculate centroids
		for j := range k.centroids {
			newCentroid := calculateCentroid(clusters[j])
			k.centroids[j] = newCentroid
		}
	}
}

func (k *KMeansCluster) Predict(data []float32) int {
	nearestCentroid := 0
	minDistance := float32(math.MaxFloat32)
	for i, centroid := range k.centroids {
		distance := euclideanDistance(data, centroid)
		if distance < minDistance {
			minDistance = distance
			nearestCentroid = i
		}
	}
	return nearestCentroid
}

func (k *KMeansCluster) Load(filename string) {
	// Code to load model from file

	json.Unmarshal(ReadFileContent(filename), &k.centroids)
}

func (k *KMeansCluster) Save(filename string) {
	JsonToFile(filename, k.centroids)
}

func euclideanDistance(a, b []float32) float32 {
	var distance float32
	for i := range a {
		distance += (a[i] - b[i]) * (a[i] - b[i])
	}
	return float32(math.Sqrt(float64(distance)))
}

func calculateCentroid(points [][]float32) []float32 {
	var centroid []float32
	if len(points) > 0 {
		// initialize centroid with zeroes
		centroid = make([]float32, len(points[0]))

		// Sum up all the points
		for _, point := range points {
			for j := range point {
				centroid[j] += point[j]
			}
		}

		// Divide by number of points to get the average
		for i := range centroid {
			centroid[i] /= float32(len(points))
		}
	}
	return centroid
}
