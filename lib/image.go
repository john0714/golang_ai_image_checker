package lib

import (
	"image"
	"path/filepath"
	"os"
	"log"
	
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
)

// private method
func cropImage(img image.Image) (image.Image, error) {
	analyzer := smartcrop.NewAnalyzer(nfnt.NewDefaultResizer())
	topCrop, err := analyzer.FindBestCrop(img, 75, 75)
	if err != nil {
		return nil, err
	}
	
	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	img = img.(SubImager).SubImage(topCrop)
	return resize.Resize(75, 75, img, resize.Lanczos3), nil
}

// public method
func DecodeImage(fname string) ([]float64, error) {
	f, err := os.Open(fname) // 파일정보를 취득해서 그 주소값을 반환함(예시:&0x000096a00)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, _, err := image.Decode(f) // 이미지를 디코딩함(인수로는 os.Open같은걸로 취득한 이미지파일 정보를 넣음) - RGB
	if err != nil {
		return nil, err
	}

	src, err = cropImage(src) // 사진을 원하는 크기로 리사이징 해서 반환(트리밍 한다고도 하고 크로핑(croping)한다고도 함) -> type image.Image
	if err != nil {
		return nil, err
	}
	
	// if you want to see resize img, do this ouput code
//	out, err := os.Create("./resize.jpg")
//    if err != nil {
//            fmt.Println(err)
//             os.Exit(1)
//    }
//	
//	var opt jpeg.Options
//  opt.Quality = 80
//  err = jpeg.Encode(out, src, &opt) // put quality to 80%
	
	bounds := src.Bounds() // Bounds returns the domain for which At can return non-zero color. -> 0이 아닌 색을 반환 할 수 있는 도메인을 반환
	w, h := bounds.Dx(), bounds.Dy() // x, y축 동일 크기로 맞춤
	if w < h {
		w = h
	} else {
		h = w
	}
	bb := make([]float64, w*h*3) // 1차원 슬라이스 생성 - rgb값을 넣기위해 3배로 지정
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := src.At(x, y).RGBA() // At returns the color of the pixel at (x, y) -> return rgba data to using RGBA()(red, green, blue, alpha = 이미지의 투명함)
			bb[y*w*3+x*3] = float64(r) / 255.0 // red값을지정 (0, 3, 6, ..., w*h*3-2)
			bb[y*w*3+x*3+1] = float64(g) / 255.0 // green값을지정 (1, 4, 7, ..., w*h*3-1)
			bb[y*w*3+x*3+2] = float64(b) / 255.0 // blue을지정 (2, 5, 8, ..., w*h*3)
		}
	}
	return bb, nil
}

func loadImageSet(c string) ([][]float64, error) {
	result := [][]float64{} // 2차원 슬라이스 생성
	f, err := os.Open(filepath.Join("dataSet", c)) // category이름(폴더이름)에 따라 폴더 오픈
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	names, err := f.Readdirnames(0) //디렉토리 안의 파일들을 읽어서 1차원 슬라이스로 넘김(값이 0보다 클경우 1개만, 0이하일경우 모든 파일을 취득)
	if err != nil {
		return nil, err
	}
	
	for _, n := range names { // range로 파일명 슬라이스를 루프(_는 key번호, n은 값 = 파일명)
		fname := filepath.Join("dataset", c, n) // 파일명과 경로를 결합해서 return함. argument에 지정한만큼 결합되므로 제일 뒤에 파일명을 적음(예시 : dataset\hamburget\9.jpg)
		log.Printf("add %q as %q", fname, c) // %q : 작은 따옴표 문자 literal(리터럴 = 소스 코드의 고정된 값)
		ff, err := DecodeImage(fname) // 지정한 각각의 파일(이미지)을 숫자 코드로 디코딩 - 1차원 슬라이스
		if err != nil {
			return nil, err
		}
		result = append(result, ff) // 2차원 슬라이스에 취득한 이미지 값(1차원 슬라이스)를 추가
	}
	
	return result, nil
}