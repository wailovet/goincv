package goincv

import (
	"image"
	"math"
	"sort"
)

type anchorsPoint struct {
	Cx     float32
	Cy     float32
	Stride float32
}

type NMS struct {
	name          string
	centerPoints  map[int][]anchorsPoint
	NMSPreMaxBox  int     //1000
	ConfThreshold float32 //0.6
	NMSThreshold  float32 //0.6
	targetWidth   int
	targetHeight  int
	sourcetWidth  int
	sourcetHeight int
	ScaleParams   ScaleParams

	boxCollection []Box
	numAnchors    int

	initFlage bool
}

func NewNMSWithName(name string, src image.Rectangle, targetWidth, targetHeight int) *NMS {
	nms := NewNMS(src, targetWidth, targetHeight)
	nms.name = name
	return nms
}
func NewNMS(src image.Rectangle, targetWidth, targetHeight int) *NMS {
	nms := &NMS{}

	nms.sourcetWidth = src.Bounds().Dx()
	nms.sourcetHeight = src.Bounds().Dy()

	var newHeight = targetHeight
	var newWidth = targetWidth

	if float64(nms.sourcetWidth)/float64(nms.sourcetHeight) > float64(targetWidth)/float64(targetHeight) {
		nms.ScaleParams.Ratio = float32(targetWidth) / float32(nms.sourcetWidth)
		newHeight = int(float32(nms.sourcetHeight) * nms.ScaleParams.Ratio)
	} else {
		nms.ScaleParams.Ratio = float32(targetHeight) / float32(nms.sourcetHeight)
		newWidth = int(float32(nms.sourcetWidth) * nms.ScaleParams.Ratio)
	}

	if newWidth < targetWidth {
		nms.ScaleParams.Dw = (targetWidth - newWidth) / 2
	}
	if newHeight < targetHeight {
		nms.ScaleParams.Dh = (targetHeight - newHeight) / 2
	}
	nms.targetWidth = targetWidth
	nms.targetHeight = targetHeight
	nms.init()
	return nms
}

func (m *NMS) init() {
	if !m.initFlage {
		if m.numAnchors == 0 {
			m.numAnchors = 2
		}
		if m.centerPoints == nil {
			m.centerPoints = map[int][]anchorsPoint{}
		}
		if m.ConfThreshold == 0 {
			m.ConfThreshold = 0.45
		}
		if m.NMSThreshold == 0 {
			m.NMSThreshold = 0.6
		}
		if m.NMSPreMaxBox == 0 {
			m.NMSPreMaxBox = 1000
		}
		m.initFlage = true
	}
}

func (m *NMS) AddValue(stride int, scoreValues []float32, bboxValues_Nx4 [][]float32, kpsValues_Nx10 [][]float32) *NMS {
	if !m.initFlage {
		m.init()
	}
	m.generatePoints(stride)
	m.generateBboxesSingleStride(stride, scoreValues, bboxValues_Nx4, kpsValues_Nx10)
	return m
}

func (m *NMS) End() []Box {
	sort.Slice(m.boxCollection, func(i, j int) bool {
		return m.boxCollection[i].Prob > m.boxCollection[j].Prob
	})

	inputBoxes := m.boxCollection
	thresh := m.NMSThreshold

	vArea := []int{}
	for i := range inputBoxes {
		vArea = append(vArea,
			(inputBoxes[i].Rectangle.Max.X-inputBoxes[i].Rectangle.Min.X+1)*(inputBoxes[i].Rectangle.Max.Y-inputBoxes[i].Rectangle.Min.Y+1),
		)
	}

	for i := 0; i < len(inputBoxes); i++ {
		for j := i + 1; j < len(inputBoxes); {
			xx1 := math.Max(float64(inputBoxes[i].Rectangle.Min.X), float64(inputBoxes[j].Rectangle.Min.X))
			yy1 := math.Max(float64(inputBoxes[i].Rectangle.Min.Y), float64(inputBoxes[j].Rectangle.Min.Y))
			xx2 := math.Min(float64(inputBoxes[i].Rectangle.Max.X), float64(inputBoxes[j].Rectangle.Max.X))
			yy2 := math.Min(float64(inputBoxes[i].Rectangle.Max.Y), float64(inputBoxes[j].Rectangle.Max.Y))
			w := math.Max(0, xx2-xx1+1)
			h := math.Max(0, yy2-yy1+1)
			inter := float32(w * h)
			ovr := inter / (float32(vArea[i]) + float32(vArea[j]) - inter)
			if ovr >= thresh {
				inputBoxes = append(inputBoxes[:j], inputBoxes[j+1:]...)
				vArea = append(vArea[:j], vArea[j+1:]...)
			} else {
				j++
			}
		}
	}
	return inputBoxes
}

func (m *NMS) generatePoints(stride int) {
	if len(m.centerPoints[stride]) > 0 {
		return
	}

	targetWidth := m.targetWidth
	targetHeight := m.targetHeight

	numGridW := targetWidth / stride
	numGridH := targetHeight / stride
	for i := 0; i < numGridH; i++ {
		for j := 0; j < numGridW; j++ {
			for k := 0; k < m.numAnchors; k++ {
				m.centerPoints[stride] = append(m.centerPoints[stride], anchorsPoint{
					Cx:     float32(j),
					Cy:     float32(i),
					Stride: float32(stride),
				})
			}
		}
	}

	// log.Println("m.centerPoints[stride]:", stride, m.centerPoints[stride])
}

func (m *NMS) generateBboxesSingleStride(stride int, scoreValues []float32, bboxValues [][]float32, kpsValues [][]float32) {
	nms_pre := m.NMSPreMaxBox
	nms_pre_ := (stride / 8) * nms_pre
	if nms_pre_ < nms_pre {
		nms_pre_ = nms_pre
	}

	numPoints := len(scoreValues)

	ratio := m.ScaleParams.Ratio
	dw := float32(m.ScaleParams.Dw)
	dh := float32(m.ScaleParams.Dh)

	stridePoints := m.centerPoints[stride]

	for i := 0; i < numPoints; i++ {
		cls_conf := scoreValues[i]
		if cls_conf < m.ConfThreshold {
			// log.Println("cls_conf:", cls_conf, stride)
			continue
		}
		point := stridePoints[i]
		cx := point.Cx
		cy := point.Cy
		s := point.Stride

		offsets := bboxValues[i]
		l := offsets[0] // left
		t := offsets[1] // top
		r := offsets[2] // right
		b := offsets[3] // bottom

		box := Box{}
		x1 := ((cx-l)*s - dw) / ratio
		y1 := ((cy-t)*s - dh) / ratio

		x2 := ((cx+r)*s - dw) / ratio // cx + r x2
		y2 := ((cy+b)*s - dh) / ratio // cy + b y2

		box.Rectangle.Min.X = int(math.Max(0, float64(x1)))
		box.Rectangle.Min.Y = int(math.Max(0, float64(y1)))
		box.Rectangle.Max.X = int(math.Min(float64(m.sourcetWidth)-1, float64(x2)))
		box.Rectangle.Max.Y = int(math.Min(float64(m.sourcetHeight)-1, float64(y2)))
		box.Prob = cls_conf
		box.Extension = map[string]interface{}{}
		box.Extension["stride"] = stride

		// landmarks
		if len(kpsValues) > 0 {
			for j := 0; j < len(kpsValues[i]); j += 2 {
				kps := BoxLandmark{}
				kps_l := kpsValues[i][j]
				kps_t := kpsValues[i][j+1]
				kps_x := ((cx+kps_l)*s - dw) / ratio
				kps_y := ((cy+kps_t)*s - dh) / ratio
				kps.X = int(math.Min(math.Max(0, float64(kps_x)), float64(m.sourcetWidth)-1))
				kps.Y = int(math.Min(math.Max(0, float64(kps_y)), float64(m.sourcetHeight)-1))
				box.Landmark = append(box.Landmark, kps)
			}
		}

		m.boxCollection = append(m.boxCollection, box)
	}

	if len(m.boxCollection) > nms_pre_ {
		sort.Slice(m.boxCollection, func(i, j int) bool {
			return m.boxCollection[i].Prob > m.boxCollection[j].Prob
		})
		m.boxCollection = m.boxCollection[:nms_pre_]
	}
}
