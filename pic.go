package goincv

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"

	"github.com/anthonynsimon/bild/segment"
	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	"github.com/spf13/cast"
)

func CopyMakeBorder(img image.Image, b, t, l, r int, c color.Color) *image.RGBA {
	rgba := image.NewRGBA(image.Rectangle{
		Min: image.ZP,
		Max: image.Pt(img.Bounds().Dx()+l+r, img.Bounds().Dy()+t+b),
	})
	newRectangle := image.Rectangle{
		Min: image.Pt(l, t),
		Max: image.Pt(img.Bounds().Dx()+l, img.Bounds().Dy()+t),
	}
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	draw.Draw(rgba, newRectangle, img, image.Point{0, 0}, draw.Src)
	return rgba
}

func File2Image(fileName string) image.Image {
	data, _ := ioutil.ReadFile(fileName)
	img, _, _ := image.Decode(bytes.NewBuffer(data))

	if img != nil {
		img = ToRGBA(img)
	}
	return img
}

func ImRead(img image.Image) (ret [][][]uint8) {
	img = ToRGBA(img) //jpg精度损失,提前转为rgba格式

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	ret = make([][][]uint8, h)
	for dh := 0; dh < h; dh++ {
		ret[dh] = make([][]uint8, w)
	}
	for dh := 0; dh < h; dh++ {
		for dw := 0; dw < w; dw++ {
			r, g, b, _ := img.At(dw, dh).RGBA()
			ret[dh][dw] = []uint8{
				uint8(b),
				uint8(g),
				uint8(r),
			}
		}
	}
	return
}

func ImRead4File(fileName string) (ret [][][]uint8) {
	data, _ := ioutil.ReadFile(fileName)
	img, _, _ := image.Decode(bytes.NewBuffer(data))
	return ImRead(img)
}

func ImReadRGB(img image.Image) (ret [][][]uint8) {
	img = ToRGBA(img) //jpg精度损失,提前转为rgba格式

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	ret = make([][][]uint8, h)
	for dh := 0; dh < h; dh++ {
		ret[dh] = make([][]uint8, w)
	}
	for dh := 0; dh < h; dh++ {
		for dw := 0; dw < w; dw++ {
			r, g, b, _ := img.At(dw, dh).RGBA()
			ret[dh][dw] = []uint8{
				uint8(r),
				uint8(g),
				uint8(b),
			}
		}
	}
	return
}

func ImReadBGR(img image.Image) (ret [][][]uint8) {
	img = ToRGBA(img) //jpg精度损失,提前转为rgba格式

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	ret = make([][][]uint8, h)
	for dh := 0; dh < h; dh++ {
		ret[dh] = make([][]uint8, w)
	}
	for dh := 0; dh < h; dh++ {
		for dw := 0; dw < w; dw++ {
			r, g, b, _ := img.At(dw, dh).RGBA()
			ret[dh][dw] = []uint8{
				uint8(b),
				uint8(g),
				uint8(r),
			}
		}
	}
	return
}

func ImageToFloat32(imgData [][][]uint8) (ret [][][]float32) {
	for i := range imgData {
		ret = append(ret, [][]float32{})
		for k := range imgData[i] {
			ret[i] = append(ret[i], []float32{})
			for j := range imgData[i][k] {
				ret[i][k] = append(ret[i][k], float32(imgData[i][k][j]))
			}
		}
	}
	return ret
}

func ImageMeanNormalizeF32RGBCHW(imgData [][][]uint8, mean, norm []float32) [][][]float32 {
	ret := [][][]float32{
		{},
		{},
		{},
	}

	rs := imgData[0]
	ret[0] = make([][]float32, len(rs))
	ret[1] = make([][]float32, len(rs))
	ret[2] = make([][]float32, len(rs))
	for i := range rs {
		ret[0][i] = make([]float32, len(rs[i]))
		ret[1][i] = make([]float32, len(rs[i]))
		ret[2][i] = make([]float32, len(rs[i]))
		for k := range rs[i] {
			r := float32(imgData[0][i][k])
			g := float32(imgData[1][i][k])
			b := float32(imgData[2][i][k])
			ret[0][i][k] = float32((r/255 - mean[0]) / norm[0])
			ret[1][i][k] = float32((g/255 - mean[1]) / norm[1])
			ret[2][i][k] = float32((b/255 - mean[2]) / norm[2])
		}
	}
	return ret
}

func ImageMeanNormalizeF32RGBHWC(imgData [][][]uint8, mean, norm []float32) [][][]float32 {
	ret := [][][]float32{}

	for i := range imgData {
		ret = append(ret, [][]float32{})
		for k := range imgData[i] {
			r := float32(imgData[i][k][0])
			g := float32(imgData[i][k][1])
			b := float32(imgData[i][k][2])

			ret[i] = append(ret[i], []float32{
				float32((r/255 - mean[0]) / norm[0]),
				float32((g/255 - mean[1]) / norm[1]),
				float32((b/255 - mean[2]) / norm[2]),
			})
		}
	}
	return ret
}

func ImReadRGBCHW(img image.Image) (ret [][][]uint8) {
	img = ToRGBA(img) //jpg精度损失,提前转为rgba格式

	ret = [][][]uint8{
		{},
		{},
		{},
	}

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	ret[0] = make([][]uint8, h)
	ret[1] = make([][]uint8, h)
	ret[2] = make([][]uint8, h)
	for dh := 0; dh < h; dh++ {
		ret[0][dh] = make([]uint8, w)
		ret[1][dh] = make([]uint8, w)
		ret[2][dh] = make([]uint8, w)
	}
	for dh := 0; dh < h; dh++ {
		for dw := 0; dw < w; dw++ {
			r, g, b, _ := img.At(dw, dh).RGBA()
			ret[0][dh][dw] = uint8(r)
			ret[1][dh][dw] = uint8(g)
			ret[2][dh][dw] = uint8(b)
		}
	}
	return
}

func Rectangle(img image.Image, r image.Rectangle, color color.Color) image.Image {
	rgba := ToRGBA(img)
	drwaRect(rgba, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, color)
	return rgba
}

func ToRGBA(img image.Image) *image.RGBA {
	if rgba, isok := img.(*image.RGBA); isok {
		return rgba
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

func drwaHLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func drwaVLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func drwaRect(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	drwaHLine(img, x1, y1, x2, col)
	drwaHLine(img, x1, y2, x2, col)
	drwaVLine(img, x1, y1, y2, col)
	drwaVLine(img, x2, y1, y2, col)
}

func SaveJPEG(img image.Image, outImgPath string) {
	buf := bytes.NewBuffer(nil)
	jpeg.Encode(buf, img, &jpeg.Options{Quality: 96})
	ioutil.WriteFile(outImgPath, buf.Bytes(), 0644)
}

func SavePNG(img image.Image, outImgPath string) {
	buf := bytes.NewBuffer(nil)
	png.Encode(buf, img)
	ioutil.WriteFile(outImgPath, buf.Bytes(), 0644)
}

func ImageClip(img image.Image, x0, y0, x1, y1 int) image.Image {
	rgba := image.NewRGBA(image.Rect(0, 0, x1-x0, y1-y0))

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{x0, y0}, draw.Src)

	return rgba
}

func RGB2Image(data [][][]float32) (image.Image, error) {
	shape := GetShapeBySlice(data)
	if shape[0] == 3 {
		//
		h := shape[1]
		w := shape[2]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		r := 0
		g := 1
		b := 2
		// a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {

				if data[r][ih][iw] < 0 {
					data[r][ih][iw] = 0
				}

				if data[r][ih][iw] > 255 {
					data[r][ih][iw] = 255
				}

				if data[g][ih][iw] < 0 {
					data[g][ih][iw] = 0
				}
				if data[g][ih][iw] > 255 {
					data[g][ih][iw] = 255
				}

				if data[b][ih][iw] < 0 {
					data[b][ih][iw] = 0
				}

				if data[b][ih][iw] > 255 {
					data[b][ih][iw] = 255
				}

				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[r][ih][iw]),
					G: uint8(data[g][ih][iw]),
					B: uint8(data[b][ih][iw]),
					A: 255,
				})
			}
		}

		return rgba, nil
	} else if shape[2] == 3 {
		h := shape[0]
		w := shape[1]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		b := 0
		g := 1
		r := 2
		// a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {

				if data[ih][iw][r] < 0 {
					data[ih][iw][r] = 0
				}
				if data[ih][iw][r] > 255 {
					data[ih][iw][r] = 255
				}

				if data[ih][iw][g] < 0 {
					data[ih][iw][g] = 0
				}
				if data[ih][iw][g] > 255 {
					data[ih][iw][g] = 255
				}

				if data[ih][iw][b] < 0 {
					data[ih][iw][b] = 0
				}
				if data[ih][iw][b] > 255 {
					data[ih][iw][b] = 255
				}
				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[ih][iw][r]),
					G: uint8(data[ih][iw][g]),
					B: uint8(data[ih][iw][b]),
					A: 255,
				})
			}
		}

		return rgba, nil
	}
	return nil, errors.New("未知格式")

}

func RGB2ImageU8(data [][][]uint8) (image.Image, error) {
	shape := GetShapeBySlice(data)
	if shape[0] == 3 {
		//
		h := shape[1]
		w := shape[2]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		r := 0
		g := 1
		b := 2
		// a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {
				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[r][ih][iw]),
					G: uint8(data[g][ih][iw]),
					B: uint8(data[b][ih][iw]),
					A: 255,
				})
				// rgba.Pix[4*ih*w+iw+r] = uint8(data[r][ih][iw])
				// rgba.Pix[4*ih*w+iw+g] = uint8(data[g][ih][iw])
				// rgba.Pix[4*ih*w+iw+b] = uint8(data[b][ih][iw])
				// rgba.Pix[4*ih*w+iw+a] = uint8(255)
			}
		}

		return rgba, nil
	} else if shape[2] == 3 {
		h := shape[0]
		w := shape[1]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		b := 2
		g := 1
		r := 0
		// a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {
				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[ih][iw][r]),
					G: uint8(data[ih][iw][g]),
					B: uint8(data[ih][iw][b]),
					A: 255,
				})
			}
		}

		return rgba, nil
	}
	return nil, errors.New("未知格式")

}

func BGR2ImageU8(data [][][]uint8) (image.Image, error) {
	shape := GetShapeBySlice(data)
	if shape[0] == 3 {
		//
		w := shape[2]
		h := shape[1]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		b := 0
		g := 1
		r := 2
		a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {
				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[r][ih][iw]),
					G: uint8(data[g][ih][iw]),
					B: uint8(data[b][ih][iw]),
					A: uint8(data[a][ih][iw]),
				})
			}
		}

		return rgba, nil
	} else if shape[2] == 3 {
		//
		h := shape[0]
		w := shape[1]
		rgba := image.NewRGBA(image.Rect(0, 0, w, h))

		b := 0
		g := 1
		r := 2
		// a := 3
		for ih := 0; ih < h; ih++ {
			for iw := 0; iw < w; iw++ {
				rgba.Set(iw, ih, color.RGBA{
					R: uint8(data[ih][iw][r]),
					G: uint8(data[ih][iw][g]),
					B: uint8(data[ih][iw][b]),
					A: 255,
				})
			}
		}

		return rgba, nil
	}
	return nil, errors.New("未知格式")

}

func BGRMat2Image(bgrData []uint8, w, h int) *image.RGBA {

	tmp1, err := Reshape(bgrData, Shape{h, w, 3}, DataTypeUInt8)
	if err != nil {
		log.Println("Reshape:", err)
		return nil
	}
	tmp2 := To3DU8(tmp1)

	img, err := BGR2ImageU8(tmp2)
	if err != nil {
		log.Println("BGR2Image:", err)
		return nil
	}
	return ToRGBA(img)
}

type BoxLandmark struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Box struct {
	Rectangle image.Rectangle
	Landmark  []BoxLandmark
	Extension map[string]interface{}
	Prob      float32
}

func DetectNms(inputBoxes []Box, thresh float32) []Box {

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

func ScrfdGenerateProposals(anchors [][]float32, scale float32, inpWidth, inpHeight, stride int, pdata_score []float32, pdata_bbox []float32, pdata_kps []float32, confThreshold float32) []Box {
	/////generate proposals
	var boxes []Box
	// vector< vector<int>> landmarks;
	// float ratioh = (float)frame.rows / newh, ratiow = (float)frame.cols / neww;

	num_grid_x := int(float32(inpWidth) / float32(stride))
	num_grid_y := int(float32(inpHeight) / float32(stride))
	pdata_bbox_index := 0
	pdata_score_index := 0
	numAnchors := len(anchors)

	for q := 0; q < numAnchors; q++ {
		anchor := anchors[q]
		anchorY := anchor[1]
		anchorW := anchor[2] - anchor[0]
		anchorH := anchor[3] - anchor[1]
		for i := 0; i < num_grid_y; i++ {
			anchorX := anchor[0]
			for j := 0; j < num_grid_x; j++ {
				if pdata_score[pdata_score_index] > confThreshold {

					dx := pdata_bbox[0+pdata_bbox_index] * float32(stride)
					dy := pdata_bbox[1+pdata_bbox_index] * float32(stride)
					dw := pdata_bbox[2+pdata_bbox_index] * float32(stride)
					dh := pdata_bbox[3+pdata_bbox_index] * float32(stride)

					cx := anchorX + anchorW*0.5
					cy := anchorY + anchorH*0.5

					x0 := cx - dx
					y0 := cy - dy
					x1 := cx + dw
					y1 := cy + dh

					boxes = append(boxes, Box{
						Rectangle: image.Rect(int(x0/scale), int(y0/scale), int(x1/scale), int(y1/scale)),
						Prob:      pdata_score[pdata_score_index],
					})

				}
				pdata_score_index++
				pdata_bbox_index += 4
			}
		}
	}

	return boxes
}

func Base64ToImage(base64Str string) image.Image {
	raw, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		log.Println("Base64ToImage error:", err)
		return nil
	}

	img, _, err := image.Decode(bytes.NewBuffer(raw))
	if err != nil {
		log.Println("Base64ToImage Decode error:", err)
		return nil
	}
	return img
}

func Image2Base64(img image.Image, types ...string) (base64Str string) {
	imgType := "png"
	if len(types) > 0 {
		imgType = types[0]
	}
	var err error
	buf := bytes.NewBuffer(nil)

	switch imgType {
	case "png":
		err = png.Encode(buf, img)
		if err != nil {
			log.Println("Image2Base64 Png Encode error:", err)
			return ""
		}
	case "jpeg":
	case "jpg":
		q := 95
		if len(types) > 1 && cast.ToInt(types[1]) > 0 {
			q = cast.ToInt(types[1])
		}
		err = jpeg.Encode(buf, img, &jpeg.Options{
			Quality: q,
		})
		if err != nil {
			log.Println("Image2Base64 Jpg Encode error:", err)
			return ""
		}
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())

}

type canvasImage struct {
	Pic  image.Image
	Rect image.Rectangle
}

func RunAndSplicingAfterCutting(img image.Image, width, height int, rollStep int, call func(item image.Image) image.Image) image.Image {

	canvasImages := []canvasImage{}
	outImages := []image.Image{}

	for x := 0; x < img.Bounds().Dx(); x += rollStep {
		for y := 0; y < img.Bounds().Dy(); y += rollStep {
			x0, y0 := x, y
			x1, y1 := x0+width, y0+height
			canvasImages = append(canvasImages, canvasImage{
				Rect: image.Rect(x0, y0, x1, y1),
				Pic:  ImageClip(img, x0, y0, x1, y1),
			})
		}
	}

	var canvas *image.RGBA
	var maskLT *image.Alpha
	var maskL *image.Alpha
	var maskT *image.Alpha
	var maskM *image.Alpha

	var bs int

	for i := range canvasImages {
		nimg := call(canvasImages[i].Pic)
		outImages = append(outImages, nimg)
		if canvas == nil {
			bs = nimg.Bounds().Dx() / canvasImages[i].Rect.Dx()

			endWidth := bs * img.Bounds().Dx()
			endHeight := bs * img.Bounds().Dy()

			canvas = image.NewRGBA(image.Rect(0, 0, endWidth, endHeight))

			maskBounds := nimg.Bounds()

			maskLT = image.NewAlpha(maskBounds)
			maskM = image.NewAlpha(maskBounds)

			maskL = image.NewAlpha(maskBounds)
			maskT = image.NewAlpha(maskBounds)

			for i := range maskLT.Pix {
				maskLT.Pix[i] = 255
				maskM.Pix[i] = 255
				maskL.Pix[i] = 255
				maskT.Pix[i] = 255
			}

			for x := 0; x < maskBounds.Dx(); x++ {
				for y := 0; y < maskBounds.Dy(); y++ {
					weightX := float64(x) / float64(maskBounds.Dx()/4)
					weightX = math.Min(weightX, 1)
					weightY := float64(y) / float64(maskBounds.Dy()/4)
					weightY = math.Min(weightY, 1)
					weightX2 := float64(maskBounds.Dx()-(1+x)) / float64(maskBounds.Dx()/4)
					weightX2 = math.Min(weightX2, 1)
					weightY2 := float64(maskBounds.Dy()-(1+y)) / float64(maskBounds.Dy()/4)
					weightY2 = math.Min(weightY2, 1)

					weight1 := weightX * weightY
					weight2 := weightX2 * weightY
					weight3 := weightX * weightY2
					weight4 := weightX2 * weightY2

					weight := math.Min(weight1, weight2)
					weight = math.Min(weight, weight3)

					weight = math.Min(weight, weight4)
					if weight < 1 {
						maskM.SetAlpha(x, y, color.Alpha{A: uint8(weight * 255)})
					}

					if weight4 < 1 {
						maskLT.SetAlpha(x, y, color.Alpha{A: uint8(weight4 * 255)})
					}

					weight = math.Min(weight3, weight4)
					if weight < 1 {
						maskT.SetAlpha(x, y, color.Alpha{A: uint8(weight * 255)})
					}

					weight = math.Min(weight2, weight4)
					if weight < 1 {
						maskL.SetAlpha(x, y, color.Alpha{A: uint8(weight * 255)})
					}

				}
			}

			// SavePNG(maskL, "maskL.png")

		}

	}

	for i := len(canvasImages) - 1; i >= 0; i-- {
		ret := outImages[i]
		mask := maskLT
		if canvasImages[i].Rect.Min.X == 0 && canvasImages[i].Rect.Min.Y == 0 {
			mask = maskLT
		} else if canvasImages[i].Rect.Min.X == 0 {
			mask = maskL
		} else if canvasImages[i].Rect.Min.Y == 0 {
			mask = maskT
		} else {
			mask = maskM
		}

		draw.DrawMask(canvas, image.Rect(
			bs*canvasImages[i].Rect.Min.X,
			bs*canvasImages[i].Rect.Min.Y,
			bs*canvasImages[i].Rect.Max.X,
			bs*canvasImages[i].Rect.Max.Y,
		), ret, image.ZP, mask, image.ZP, draw.Over)
	}
	return canvas
}

func MultiImageFusion(imgs []image.Image, mode int) image.Image {
	maxX := 0
	maxY := 0

	imgsuint8 := [][][][]uint8{}
	for i := range imgs {
		if imgs[i].Bounds().Dx() > maxX {
			maxX = imgs[i].Bounds().Dx()
		}
		if imgs[i].Bounds().Dy() > maxY {
			maxY = imgs[i].Bounds().Dy()
		}
		imgsuint8 = append(imgsuint8, ImReadRGB(imgs[i]))
		// log.Println(i)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, maxX, maxY))

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			rall := []float64{}
			gall := []float64{}
			ball := []float64{}
			for i := range imgs {
				rall = append(rall, float64(imgsuint8[i][y][x][0]))
				gall = append(gall, float64(imgsuint8[i][y][x][1]))
				ball = append(ball, float64(imgsuint8[i][y][x][2]))
			}
			var rmean, gmean, bmean float64

			if mode == 0 {
				rmean = Mean(rall)
				gmean = Mean(gall)
				bmean = Mean(ball)
			}
			if mode == 1 {
				rmean = Variance(rall)
				gmean = Variance(gall)
				bmean = Variance(ball)
			}
			if mode == 2 {
				rmean = Std(rall)
				gmean = Std(gall)
				bmean = Std(ball)
			}

			rgba.Set(x, y, color.RGBA{
				R: uint8(rmean),
				G: uint8(gmean),
				B: uint8(bmean),
				A: 255,
			})

		}
	}
	return rgba

}

func BWMask2AMask(img image.Image) *image.Alpha {
	ret := image.NewAlpha(img.Bounds())
	result := segment.Threshold(img, 128)

	for i := range result.Pix {
		ret.Pix[i] = result.Pix[i]
	}
	return ret
}

func BWMask2AMaskLadder(img image.Image) *image.Alpha {
	imgSet := image.NewGray(img.Bounds())
	ret := image.NewAlpha(img.Bounds())

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			imgSet.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	for i := range imgSet.Pix {
		ret.Pix[i] = imgSet.Pix[i]
	}

	return ret
}

func ExtractionBasisMask(img image.Image, mask image.Image) image.Image {
	var canvas *image.RGBA = image.NewRGBA(img.Bounds())

	a, ok := mask.(*image.Alpha)
	if ok {
		draw.DrawMask(canvas, canvas.Bounds(), img, image.ZP, a, image.ZP, draw.Over)
		return canvas
	} else {
		draw.DrawMask(canvas, canvas.Bounds(), img, image.ZP, BWMask2AMask(mask), image.ZP, draw.Over)
		return canvas
	}
}

type ScaleParams struct {
	Ratio float32
	Dw    int
	Dh    int
}

var ResizeMode imaging.ResampleFilter = imaging.Lanczos

func ResizeImageNoDeformation(srcImage image.Image, width, height int) (ret image.Image, ratio float32) {
	ox := srcImage.Bounds().Dx()
	oy := srcImage.Bounds().Dy()

	if float64(ox)/float64(oy) > float64(width)/float64(height) {
		ratio = float32(width) / float32(ox)
		height = int(float32(oy) * ratio)
	} else {
		ratio = float32(height) / float32(oy)
		width = int(float32(ox) * ratio)
	}

	ret = imaging.Resize(srcImage, width, height, ResizeMode)
	return
}

func ResizeImageBorder(srcImage image.Image, width, height int, bg color.Color) (ret image.Image, scaleParams ScaleParams) {

	ret, scaleParams.Ratio = ResizeImageNoDeformation(srcImage, width, height)
	// log.Println("scaleParams.Ratio:", scaleParams.Ratio)

	if ret.Bounds().Dx() < width {
		scaleParams.Dw = (width - ret.Bounds().Dx()) / 2

	}
	if ret.Bounds().Dy() < height {
		scaleParams.Dh = (height - ret.Bounds().Dy()) / 2
	}

	ret = CopyMakeBorder(ret, scaleParams.Dh, scaleParams.Dh, scaleParams.Dw, scaleParams.Dw, bg)
	if ret.Bounds().Dx() != width || ret.Bounds().Dy() != height {
		ret = imaging.Resize(ret, width, height, ResizeMode)
	}
	return
}

func MergeAlpha(ia []*image.Alpha) *image.Alpha {
	c := TimeConsuming("MergeAlpha")
	defer c()
	ret := image.NewAlpha(ia[0].Rect)
	for x := 0; x < ret.Rect.Max.X; x++ {
		for y := 0; y < ret.Rect.Max.Y; y++ {
			c := 0
			for i := range ia {
				if ia[i] != nil {
					c += int(ia[i].AlphaAt(x, y).A)
				}
			}
			if c > 255 {
				c = 255
			}
			ret.SetAlpha(x, y, color.Alpha{
				A: uint8(c),
			})
		}
	}
	// for k := range ret.Pix {
	// 	if len(ret.Pix) <= k {
	// 		ret.Pix = append(ret.Pix, ia[0].Pix[k])
	// 	}
	// 	for i := range ia {
	// 		if ret.Pix[k]+ia[i].Pix[k] < 255 {
	// 			ret.Pix[k] += ia[i].Pix[k]
	// 		} else {
	// 			ret.Pix[k] = 255
	// 		}
	// 	}
	// }
	return ret
}

func CompareImage(a, b image.Image, w, h int) float32 {
	hash1, _ := goimagehash.ExtDifferenceHash(a, w, h)
	hash2, _ := goimagehash.ExtDifferenceHash(b, w, h)
	distance, _ := hash1.Distance(hash2)
	return float32(distance) / float32(w*h)
}
