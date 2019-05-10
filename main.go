package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
)

func main() {
	// 清理image 目录
	dirList, err := ioutil.ReadDir("./images/")
	if err != nil {
		fmt.Println("read images dir error")
	}
	for _, v := range dirList {
		if v.Name() != ".gitignore" {
			_ = os.Remove("./images/" + v.Name())
		}
	}
	//创建一个随机地图；
	randomMap := NewPointMap(100)
	//设置图像的内容与地图大小一致；
	img := image.NewNRGBA(image.Rect(0, 0, randomMap.Size, randomMap.Size))
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
	aStar.SaveImage(img)
	aStar.RunAndSaveImage(img)
}
