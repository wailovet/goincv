package goincv

import (
	"math"
	"math/rand"
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
}

func (k *KMeansCluster) Save(filename string) {
	// Code to save model to file
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

type DBScanCluster struct {
	data    [][]float32 // 数据
	eps     float32     // eps为领域半径
	minPts  int         // minPts为密度阈值
	labels  []int       // 每个点所属类别
	visited []bool      // 每个点是否已经被访问过
	cluster int         // 当前类别
}

func (d *DBScanCluster) Add(data []float32) {
	d.data = append(d.data, data)
}

func (d *DBScanCluster) Learn(eps float32, minPts int) {
	d.eps = eps
	d.minPts = minPts
	d.labels = make([]int, len(d.data))
	d.visited = make([]bool, len(d.data))
	d.cluster = 0

	// 初始化类别为-1（表示未被访问过）
	for i := range d.labels {
		d.labels[i] = -1
	}

	// 遍历每个点
	for i := range d.data {
		if !d.visited[i] {
			d.visited[i] = true

			// 找到点i的邻域内的点
			neighbours := d.regionQuery(i)

			// 如果邻域内点数小于minPts，则标记为噪点
			if len(neighbours) < d.minPts {
				d.labels[i] = -1
			} else {
				// 否则，开始新的一类
				d.cluster++
				d.expandCluster(i, neighbours)
			}
		}
	}
}

func (d *DBScanCluster) expandCluster(point int, neighbours []int) {
	d.labels[point] = d.cluster

	for _, neighbour := range neighbours {
		if !d.visited[neighbour] {
			d.visited[neighbour] = true
			newNeighbours := d.regionQuery(neighbour)
			if len(newNeighbours) >= d.minPts {
				neighbours = append(neighbours, newNeighbours...)
			}
		}
		if d.labels[neighbour] == -1 {
			d.labels[neighbour] = d.cluster
		}
	}
}

func (d *DBScanCluster) regionQuery(point int) []int {
	var neighbours []int
	for i, data := range d.data {
		if euclideanDistance(d.data[point], data) <= d.eps {
			neighbours = append(neighbours, i)
		}
	}
	return neighbours
}

func (d *DBScanCluster) Predict(data []float32) int {
	for i, point := range d.data {
		if euclideanDistance(point, data) == 0 {
			return d.labels[i]
		}
	}
	return -1
}

func (d *DBScanCluster) Load(filename string) {
	// 从文件中读取数据，并反序列化
	// ...
}

func (d *DBScanCluster) Save(filename string) {
	// 序列化并保存数据到文件
	// ...
}

type AgglomerativeCluster struct {
	data      [][]float32 // 数据
	linkage   string      // linkage 用来指定聚类时的距离度量方式
	labels    []int       // 每个点的类别标签
	nClusters int         //类别数量
}

func (a *AgglomerativeCluster) Add(data []float32) {
	a.data = append(a.data, data)
}

func (a *AgglomerativeCluster) Learn(linkage string, nClusters int) {
	a.linkage = linkage
	a.nClusters = nClusters
	nPoints := len(a.data)
	a.labels = make([]int, nPoints)
	// 初始化标签为点的索引
	for i := range a.labels {
		a.labels[i] = i
	}
	// 逐步合并类别
	for nClusters > a.nClusters {
		var minDist float32
		var minI, minJ int
		for i := 0; i < nPoints; i++ {
			for j := i + 1; j < nPoints; j++ {
				dist := distance(a.data[i], a.data[j], a.linkage)
				if a.labels[i] != a.labels[j] && (minI == 0 && minJ == 0 || dist < minDist) {
					minDist = dist
					minI = i
					minJ = j
				}
			}
		}
		// 将类别minJ的所有点的标签更新为minI
		for i := 0; i < nPoints; i++ {
			if a.labels[i] == a.labels[minJ] {
				a.labels[i] = a.labels[minI]
			}
		}
		nClusters--
	}
}

func (a *AgglomerativeCluster) Predict(data []float32) int {
	var minDist float32
	var minI int
	for i, point := range a.data {
		dist := euclideanDistance(point, data)
		if minI == 0 || dist < minDist {
			minDist = dist
			minI = i
		}
	}
	return a.labels[minI]
}

func (a *AgglomerativeCluster) Load(filename string) {
	// 从文件中读取数据，并反序列化
	// ...
}

func (a *AgglomerativeCluster) Save(filename string) {
	// 序列化并保存数据到文件
	// ...
}

// distance 计算距离
func distance(p1, p2 []float32, linkage string) float32 {
	switch linkage {
	case "single":
		return singleLinkage(p1, p2)
	case "complete":
		return completeLinkage(p1, p2)
	case "average":
		return averageLinkage(p1, p2)
	default:
		return euclideanDistance(p1, p2)
	}
}

// singleLinkage 计算单链接距离
func singleLinkage(p1, p2 []float32) float32 {
	var minDist float32
	for i := range p1 {
		dist := math.Abs(float64(p1[i] - p2[i]))
		if i == 0 || dist < float64(minDist) {
			minDist = float32(dist)
		}
	}
	return minDist
}

// completeLinkage 计算完全链接距离
func completeLinkage(p1, p2 []float32) float32 {
	var maxDist float32
	for i := range p1 {
		dist := math.Abs(float64(p1[i] - p2[i]))
		if dist > float64(maxDist) {
			maxDist = float32(dist)
		}
	}
	return maxDist
}

// averageLinkage 计算平均链接距离
func averageLinkage(p1, p2 []float32) float32 {
	var sumDist float32
	for i := range p1 {
		dist := math.Abs(float64(p1[i] - p2[i]))
		sumDist += float32(dist)
	}
	return sumDist / float32(len(p1))
}

type KMeansPlusPlus struct {
	data           [][]float32 // 数据
	nClusters      int         // 类别数量
	clusterCenters [][]float32 //类中心
	labels         []int       // 每个点的类别标签
}

func (k *KMeansPlusPlus) Add(data []float32) {
	k.data = append(k.data, data)
}

func (k *KMeansPlusPlus) Learn(nClusters int, iterations int) {
	k.nClusters = nClusters
	k.clusterCenters = k.initClusters()
	k.labels = make([]int, len(k.data))

	var clusterChanged bool
	var iteration int
	for {
		clusterChanged = k.assignLabels()
		if !clusterChanged || iteration > iterations {
			break
		}
		k.updateClusters()
		iteration++
	}
}

func (k *KMeansPlusPlus) initClusters() [][]float32 {
	// 随机选取一个点作为第一个类中心
	firstCenter := rand.Intn(len(k.data))
	clusterCenters := [][]float32{k.data[firstCenter]}

	for i := 1; i < k.nClusters; i++ {
		var distances []float32
		for _, point := range k.data {
			var minDistance float32 = math.MaxFloat32
			for _, center := range clusterCenters {
				dist := euclideanDistance(point, center)
				if dist < minDistance {
					minDistance = dist
				}
			}
			distances = append(distances, minDistance)
		}
		newClusterCenter := k.weightedRandom(distances)
		clusterCenters = append(clusterCenters, k.data[newClusterCenter])
	}
	return clusterCenters
}

// weightedRandom 通过给定的概率分布随机选取一个点
func (k *KMeansPlusPlus) weightedRandom(weights []float32) int {
	var sumWeights float32
	for _, weight := range weights {
		sumWeights += weight
	}
	randValue := rand.Float32() * sumWeights
	for i, weight := range weights {
		randValue -= weight
		if randValue < 0 {
			return i
		}
	}
	return len(weights) - 1
}
func (k *KMeansPlusPlus) assignLabels() bool {
	var clusterChanged bool
	for i, point := range k.data {
		var minDistance float32 = math.MaxFloat32
		minCluster := 0
		for j, center := range k.clusterCenters {
			dist := euclideanDistance(point, center)
			if dist < minDistance {
				minDistance = dist
				minCluster = j
			}
		}
		if k.labels[i] != minCluster {
			clusterChanged = true
		}
		k.labels[i] = minCluster
	}
	return clusterChanged
}

func (k *KMeansPlusPlus) updateClusters() {
	for i := 0; i < k.nClusters; i++ {
		var clusterPoints [][]float32
		for j, label := range k.labels {
			if label == i {
				clusterPoints = append(clusterPoints, k.data[j])
			}
		}
		k.clusterCenters[i] = calculateCentroid(clusterPoints)
	}
}

func (k *KMeansPlusPlus) Predict(data []float32) int {
	var minDistance float32 = math.MaxFloat32
	minCluster := 0
	for i, center := range k.clusterCenters {
		dist := euclideanDistance(data, center)
		if dist < minDistance {
			minDistance = dist
			minCluster = i
		}
	}
	return minCluster
}
