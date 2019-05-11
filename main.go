package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"os"
)

var output string
var outputGif bool
var outputImg bool

func init() {
	flag.BoolVar(&outputImg, "build_image", false, "是否生成图片")
	flag.BoolVar(&outputGif, "build_gif", false, "是否生成GIF文件")
	flag.StringVar(&output, "o", "output.gif", "生成gif的文件名")
	flag.Parse()
	// 清理image 目录
	fmt.Printf("===== clean './images/' directory \n")
	dirList, err := ioutil.ReadDir("./images/")
	if err != nil {
		fmt.Println("read images dir error")
		return
	}
	for _, v := range dirList {
		if v.Name() != ".gitignore" {
			_ = os.Remove("./images/" + v.Name())
		}
	}
}

func main() {
	var f *os.File
	var err error
	if outputGif {
		f, err = os.Create(fmt.Sprintf("%v", output))
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		defer f.Close()
	}

	//创建一个随机地图；
	randomMap := NewPointMap(100)
	//设置图像的内容与地图大小一致；
	img := image.NewNRGBA(image.Rect(0, 0, randomMap.Size, randomMap.Size))
	anim := gif.GIF{}
	x := make([]int, randomMap.Size)
	//绘制地图：对于障碍物绘制一个灰色的方块，其他区域绘制一个白色的的方块；
	for i := range x {
		for j := range x {
			if randomMap.IsObstacle(i, j) {
				// 绘制障碍物
				img.Set(i, j, color.RGBA{R: uint8(128), G: uint8(128), B: uint8(128), A: uint8(255)})
			} else {
				// 绘制空地
				img.Set(i, j, color.RGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)})
			}
		}
	}
	//绘制起点为蓝色方块；
	img.Set(0, 0, color.RGBA{R: uint8(0), G: uint8(0), B: uint8(255), A: uint8(255)})
	//绘制终点为红色方块；
	img.Set(randomMap.Size-1, randomMap.Size-1, color.RGBA{R: uint8(255), G: uint8(0), B: uint8(0), A: uint8(255)})
	//调用算法来查找路径；
	aStar := NewAStar(randomMap)
	if outputImg {
		aStar.SaveImage(img)
	}
	if outputGif && !outputImg {
		aStar.MakeGifFrame(img, &anim)
	}
	aStar.RunAndSaveImage(img, &anim, outputImg, outputGif)
	//制作GIF
	if outputGif && outputImg {
		aStar.BuildGif(output)
	}
	if outputGif && !outputImg {
		_ = gif.EncodeAll(f, &anim)
	}

}
