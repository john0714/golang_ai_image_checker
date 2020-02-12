package lib

import (
	"bufio"
	"errors"
	"os"
	"encoding/json"
	"fmt"
	
	"github.com/goml/gobrain"
)

// private method
func bin(n int) []float64 {
	f := [8]float64{} // MAX model 8개
	for i := uint(0); i < 8; i++ {
		f[i] = float64((n >> i) & 1) // shift로 이진수 형태로 저장
	}
	return f[:] // 전체 내용을 출력하기위해 [:]사용(기존 형태는 [8]로 지정된 형태이므로 []float64형태의 값으로 return 불가)
}

// public method(just a function with receiver argument)
func LoadModel() (*gobrain.FeedForward, []string, error) {
	f, err := os.Open("labels.txt")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	labels := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, nil, err
	}

	if len(labels) == 0 {
		return nil, nil, errors.New("No labels found")
	}

	// brain = 모델
	f, err = os.Open("brain.json")
	if err != nil {
		return nil, labels, nil
	}
	defer f.Close()

	ff := &gobrain.FeedForward{}
	err = json.NewDecoder(f).Decode(ff)
	if err != nil {
		return nil, labels, err
	}
	return ff, labels, nil
}

// public method(just a function with receiver argument)
func MakeModel(labels []string) (*gobrain.FeedForward, error) {
	ff := &gobrain.FeedForward{}
	patterns := [][][]float64{} // 3차원 모델패턴 slice
	
	for i, c := range labels { //Label slice나누기(i = 번호, category = 값)
		bset, err := loadImageSet(c) // 이미지를 2차원 슳라이스로 리턴 - 이미지 숫자만큼 값이 지정됨
		if err != nil {
			return nil, err
		}
		
		//  모든 이미지가 패턴에 {{이미지 값}{라벨 값}}의 형태로 저장됨	
		for _, b := range bset {
			patterns = append(patterns, [][]float64{b, bin(i)}) // 3차원 패턴 슬라이드에 2차원 이미지 슬라이드를 추가 , 저장할 모델의 위치값를 binary형태로 지정(3차원 슬라이스로 전부 지정됨)
		}
		// fmt.Println(bin(i)) // [0 0 0 0 0 0 0 0], [1 0 0 0 0 0 0 0], [0 1 0 0 0 0 0 0]... 이런식으로 각 이미지의 위치가 지정됨. 결과에선 각 위치에따른 일치 퍼센트가 나옴
	}
	
	if len(patterns) == 0 || len(patterns[0][0]) == 0 {
		return nil, errors.New("No images found")
	}
	
	fmt.Println("training now... please wait...")
	ff.Init(len(patterns[0][0]), 40, len(patterns[0][1])) // input(입력) = 이미지 값, hidden(기억) = 기억노드 수 지정, output(출력) = 라벨값
	ff.Train(patterns, 1000, 0.6, 0.4, false) // 학습(패턴, 학습횟수, 상수, 상수, 에러값 출력 플래그)
	return ff, nil
}

// public method
func SaveModel(ff *gobrain.FeedForward) error {
	f, err := os.Create("brain.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(ff)
}