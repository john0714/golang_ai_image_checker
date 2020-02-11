package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"goImageChecker/lib"
)

func check(d []float64) int {
	n := 0
	for i, v := range d {
		if v > 0.9 { // 이미지의 일치값이 90%이상일 경우 맞다고 판단
			n += 1 << uint(i)
		}
	}
	return n
}

func main() {
	// Argument파싱
	flag.Parse()

	// AI 모델 취득 -> brain.json은 데이터셋을 보고 자동생성됨
	ff, labels, err := lib.LoadModel()
	if err != nil {
		log.Fatal(err)
	}
	
	// 모델이 없으면 이미지를 기반으로 모델을 생성합니다.
	if ff == nil {
		log.Println("making model file since not found")
		ff, err = lib.MakeModel(labels) // 학습된 모델을 생성
		if err != nil {
			log.Fatal(err)
		}
		err = lib.SaveModel(ff) // 학습된 모델을저장(json)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 인수 존재체크
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// 입력받은 이미지 판단
	for _, arg := range flag.Args() {
		input, err := lib.DecodeImage(arg)
		if err != nil {
			log.Fatal(err)
		}
		n := check(ff.Update(input)) // AI판단값 체크
		if n >= 0 && n < len(labels) { // 일치한다고 판단되는 종류가 있다면 그 종류를 출력
			fmt.Println(labels[n])
		} else { // 판단 된게 없다면 종류가 없다고 출력
			fmt.Println("unknown image")
		}
	}
}