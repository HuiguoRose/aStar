package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"os"
	"time"
)

type AStar struct {
	PointMap  *PointMap
	OpenSet   []*Point
	CloseSet  []*Point
	imageList []string
}

func NewAStar(pointMap *PointMap) *AStar {
	return &AStar{PointMap: pointMap}
}

//节点到起点的移动代价，对应了上文的g(n)。
func (a *AStar) BaseCost(p *Point) float64 {
	xDis := float64(p.X)
	yDis := float64(p.Y)
	// Distance to start point
	return xDis + yDis + (math.Sqrt(2)-2)*math.Min(xDis, yDis)
}

//节点到终点的启发函数，对应上文的h(n)。由于我们是基于网格的图形，所以这个函数和上一个函数用的是对角距离。
func (a *AStar) HeuristicCost(p *Point) float64 {
	xDis := float64(a.PointMap.Size - 1 - p.X)
	yDis := float64(a.PointMap.Size - 1 - p.Y)
	// Distance to end point
	return xDis + yDis + (math.Sqrt(2)-2)*math.Min(xDis, yDis)
}

//代价总和，即对应上面提到的f(n)。
func (a *AStar) TotalCost(p *Point) int {
	return int(a.BaseCost(p) + a.HeuristicCost(p))
}

//判断点是否有效，不在地图内部或者障碍物所在点都是无效的。
func (a *AStar) IsValidPoint(x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if x >= a.PointMap.Size || y >= a.PointMap.Size {
		return false
	}
	return !a.PointMap.IsObstacle(x, y)
}

//判断点是否在某个集合中。
func (a *AStar) IsInPointList(p *Point, pointList []*Point) bool {
	for _, point := range pointList {
		if point.X == p.X && point.Y == p.Y {
			return true
		}
	}
	return false
}

//判断点是否在open_set中。
func (a *AStar) IsInOpenList(p *Point) bool {
	return a.IsInPointList(p, a.OpenSet)
}

//判断点是否在close_set中。
func (a *AStar) IsInCloseList(p *Point) bool {
	return a.IsInPointList(p, a.CloseSet)
}

//判断点是否是起点。
func (a *AStar) IsStartPoint(p *Point) bool {
	return p.X == 0 && p.Y == 0
}

//判断点是否是终点。
func (a *AStar) IsEndPoint(p *Point) bool {
	return p.X == a.PointMap.Size-1 && p.Y == a.PointMap.Size-1
}

func (a *AStar) RunAndSaveImage(img *image.NRGBA) {
	startTime := time.Now().UnixNano()
	startPoint := NewPoint(0, 0) //起点
	startPoint.Cost = 0
	a.OpenSet = append(a.OpenSet, startPoint) //起点添加到开启集合
	for {
		index := a.SelectPointInOpenList() //选择最优点
		if index < 0 {
			fmt.Println("No path found, algorithm failed!!!")
			return
		}
		p := a.OpenSet[index]
		img.Set(p.X, p.Y, color.RGBA{R: uint8(0), G: uint8(191), B: uint8(191), A: uint8(255)})
		a.SaveImage(img)
		if a.IsEndPoint(p) { //如果是终点画出所有最优点的坐标
			a.BuildPath(p, img, startTime)
			return
		}
		a.OpenSet = DeletePointSliceIndex(a.OpenSet, index)
		//a.OpenSet= append(a.OpenSet[:index], a.OpenSet[index+1:]...)
		a.CloseSet = append(a.CloseSet, p)
		//# Process all neighbors
		x := p.X
		y := p.Y
		a.ProcessPoint(x-1, y+1, p)
		a.ProcessPoint(x-1, y, p)
		a.ProcessPoint(x-1, y-1, p)
		a.ProcessPoint(x, y-1, p)
		a.ProcessPoint(x+1, y-1, p)
		a.ProcessPoint(x+1, y, p)
		a.ProcessPoint(x+1, y+1, p)
		a.ProcessPoint(x, y+1, p)
	}
}

//从open_set中找到优先级最高的节点，返回其索引。
func (a *AStar) SelectPointInOpenList() int {
	index := 0
	selectIndex := -1
	minCost := math.MaxInt64
	for _, p := range a.OpenSet {
		cost := a.TotalCost(p)
		if cost < minCost {
			minCost = cost
			selectIndex = index
		}
		index += 1
	}
	return selectIndex
}

//从终点往回沿着parent构造结果路径。然后从起点开始绘制结果，结果使用绿色方块，每次绘制一步便保存一个图片。
func (a *AStar) BuildPath(p *Point, img *image.NRGBA, startTime int64) {
	var path []*Point
	for {
		path = InsertPointSliceCopy(path, []*Point{p}, 0)
		if a.IsStartPoint(p) {
			break
		} else {
			p = p.Parent
		}
	}
	for _, p := range path {
		img.Set(p.X, p.Y, color.RGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)})
		a.SaveImage(img)
	}
	endTime := time.Now().UnixNano()
	fmt.Printf("===== Algorithm finish in %d ms\n", int(endTime-startTime)/1e6)
}

//将当前状态保存到图片中，图片以当前时间命名。
func (a *AStar) SaveImage(img *image.NRGBA) {
	filename := fmt.Sprintf("./images/%v.png", time.Now().UnixNano())
	a.imageList = append(a.imageList, filename)
	imgFile, _ := os.Create(filename)
	defer imgFile.Close()
	err := png.Encode(imgFile, img)
	if err != nil {
		fmt.Println(err)
	}
}

//针对每一个节点进行处理：如果是没有处理过的节点，则计算优先级设置父节点，并且添加到open_set中。
func (a *AStar) ProcessPoint(x, y int, parent *Point) {
	if !a.IsValidPoint(x, y) {
		//Do nothing for invalid point
		return
	}
	p := NewPoint(x, y)
	if a.IsInCloseList(p) {
		//Do nothing for visited point
		return
	}
	if !a.IsInOpenList(p) {
		p.Parent = parent
		p.Cost = a.TotalCost(p)
		a.OpenSet = append(a.OpenSet, p) //将临点添加到open_set
	}
}

func (a *AStar) BuildGif(fileName string) {
	fmt.Printf("===== building gif file Please wait\n")
	anim := gif.GIF{}
	for _, file := range a.imageList {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Could not open file %s. Error: %s\n", file, err)
			return
		}
		img, _, _ := image.Decode(f)
		paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)

		anim.Image = append(anim.Image, paletted)
		anim.Delay = append(anim.Delay, 0)
	}
	f, _ := os.Create(fmt.Sprintf("%v.gif", fileName))
	defer f.Close()
	_ = gif.EncodeAll(f, &anim)
	fmt.Printf("===== building gif file success\n")
}
