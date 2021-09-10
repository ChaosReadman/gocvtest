package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/yofu/dxf"
	"gocv.io/x/gocv"
)

func writeImg(file string, img *gocv.Mat) {
	if ok := gocv.IMWrite(file, *img); !ok {
		fmt.Printf("Failed to write image: %s\n", file)
		os.Exit(1)
	}
}

func readImg(file string) gocv.Mat {
	img := gocv.IMRead(file, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("Could not read image %s\n", file)
		os.Exit(1)
	}
	return img
}

func DrawContourToDXF(name string, c gocv.PointsVector) {
	d := dxf.NewDrawing()
	d.Header().LtScale = 100.0
	d.AddLayer("flat", dxf.DefaultColor, dxf.DefaultLineType, true)
	d.ChangeLayer("flat")

	for i := 0; i < c.Size(); i++ {
		for j := 0; j < c.At(i).Size()-1; j++ {
			x1 := float64(c.At(i).At(j).X)
			y1 := float64(c.At(i).At(j).Y)
			x2 := float64(c.At(i).At(j + 1).X)
			y2 := float64(c.At(i).At(j + 1).Y)
			d.Line(x1, -y1, 0, x2, -y2, 0)
		}
	}

	err := d.SaveAs(name)
	if err != nil {
		return
	}
}

func main() {

	//コマンドライン引数準備
	args := os.Args[1]

	//読み込み画像拡張子準備
	jpg := ".jpg"

	//出力ファイル名準備
	gra := "_gray.jpg"
	thr := "_Threshold.jpg"
	con := "_Contour.jpg"
	tes := "_test.dxf"

	imgName := args + jpg
	imgGrayName := args + gra
	imgTresholdName := args + thr
	imgContourName := args + con
	testDxf := args + tes

	fmt.Println(args)

	// ファイル読み込み
	img := readImg(imgName)

	// グレースケール作成
	imgGray := gocv.NewMat()
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)
	// ここまでを保存
	writeImg(imgGrayName, &imgGray)

	imgTreshold := gocv.NewMat()
	gocv.Threshold(imgGray, &imgTreshold, 127, 255, 0)
	// ここまでを保存
	writeImg(imgTresholdName, &imgTreshold)

	// 輪郭を抽出
	contours := gocv.FindContours(imgTreshold, gocv.RetrievalTree, gocv.ChainApproxSimple)
	// 元のイメージに書き込む
	gocv.DrawContours(&img, contours, -1, color.RGBA{255, 0, 0, 0}, 3)
	// ここまでを保存
	writeImg(imgContourName, &img)

	// DXFで出力
	DrawContourToDXF(testDxf, contours)
}
