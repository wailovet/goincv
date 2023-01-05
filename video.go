package goincv

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
)

var FFmpegPath = "ffmpeg"

func FFmpegCap(videoPath string, picFormat string, r int, temp string, extArgs []string) ([]string, error) {

	args := []string{"-i", videoPath, "-r", fmt.Sprint(r), "-f", "image2"}
	args = append(args, extArgs...)
	args = append(args, filepath.Join(temp, "%05d."+picFormat))
	log.Println("args:", args)
	_, err := exec.Command(FFmpegPath, args...).CombinedOutput()
	if err != nil {
		return nil, err
	}
	absTemp, _ := filepath.Abs(temp)
	fs, err := os.ReadDir(absTemp)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for i := range fs {
		ret = append(ret, filepath.Join(absTemp, fs[i].Name()))
	}
	return ret, nil
}

func FFmpegWav(videoPath string, wav string) error {
	_, err := exec.Command(FFmpegPath, "-i", videoPath, "-f", "wav", wav).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func FFmpegSynthesis(outVideoPath string, picFormat string, r int, temp string, wavFile string) error {

	return WriteFileWithURICallback(outVideoPath, func(tmpFile string) error {
		args := []string{"-r", fmt.Sprint(r), "-i", filepath.Join(temp, "%05d."+picFormat), "-i", wavFile, "-c:a", "aac", "-strict", "experimental", "-vcodec", "libx264", "-pix_fmt", "yuv420p", "-y", tmpFile}
		if wavFile == "" {
			args = []string{"-r", fmt.Sprint(r), "-i", filepath.Join(temp, "%05d."+picFormat), "-vcodec", "libx264", "-pix_fmt", "yuv420p", "-y", tmpFile}
		}
		log.Println("FFmpegSynthesis:", args)
		data, err := exec.Command(FFmpegPath, args...).CombinedOutput()
		if err != nil {
			return errors.New(string(data))
		}
		return nil
	})

}

func FFmpegWorking(inVideoPath string, outVideoPath string, picFormat string, r int, temp string, ext []string, do func(imgPath string) error) error {
	log.Println("解码视频中...")
	data, err := FFmpegCap(inVideoPath, picFormat, r, temp, ext)
	if err != nil {
		return err
	}
	tmpWavFile := filepath.Join(os.TempDir(), "tmp.wav")
	FFmpegWav(inVideoPath, tmpWavFile)
	defer os.RemoveAll(tmpWavFile)

	log.Println("解码完成...")
	for i := range data {
		err := do(data[i])
		if err != nil {
			return err
		}
	}
	return FFmpegSynthesis(outVideoPath, picFormat, r, temp, tmpWavFile)
}

func FFmpegCapOnec(videoPath string, s int) (string, error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprint(time.Now().Format("20060102150405"), "_", time.Now().UnixNano(), "_", rand.Intn(999999), ".jpg"))
	args := []string{"-i", videoPath, "-y", "-ss", fmt.Sprint(s), tmp}

	data, _ := exec.Command(FFmpegPath, args...).CombinedOutput()
	picRaw, _ := ioutil.ReadFile(tmp)
	if len(picRaw) == 0 {
		log.Println("FFmpegCapOnec :", string(data))
		return "", fmt.Errorf("图片未生成")
	}

	return tmp, nil
}

func FFmpegMultiImageFusion(videoPath string, width, height, r int, temp string) (image.Image, error) {
	data, err := FFmpegCap(videoPath, "png", r, temp, []string{
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
	})
	if err != nil {
		return nil, err
	}

	imgs := []image.Image{}
	for i := range data {
		item := File2Image(data[i])
		imgs = append(imgs, item)
	}
	img := MultiImageFusion(imgs, 1)
	img = effect.Sobel(img)
	img = segment.Threshold(img, 128)
	return img, nil
}

func FFmpegMultiImageFusion4s(videoPath string, width, height int, temp string) (ret []image.Image, err error) {
	r := 8
	secInterval := 1

	data, err := FFmpegCap(videoPath, "png", r, temp, []string{
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
	})
	if err != nil {
		return nil, err
	}

	imgs := []image.Image{}
	for i := range data {
		item := File2Image(data[i])
		imgs = append(imgs, item)
		if len(imgs) >= r*2 {

			img := MultiImageFusion(imgs, 1)
			img = effect.Sobel(img)
			img = segment.Threshold(img, 128)

			canvas := ToRGBA(imgs[0])
			draw.DrawMask(canvas, canvas.Bounds(), img, image.ZP, BWMask2AMask(img), image.ZP, draw.Over)

			ret = append(ret, canvas)
			imgs = imgs[r*secInterval:]
		}
	}

	img := MultiImageFusion(imgs, 1)

	img = effect.Sobel(img)
	img = segment.Threshold(img, 128)

	canvas := ToRGBA(imgs[0])
	draw.DrawMask(canvas, canvas.Bounds(), img, image.ZP, BWMask2AMask(img), image.ZP, draw.Over)

	ret = append(ret, canvas)

	return ret, nil
}
