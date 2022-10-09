package Processing

import (
	"log"
	"regexp"
)

/*
//총 모험섬 종류 및 시간
type Island struct {
	Island     []IslandElement `json:"Island"`
	IslandDate string          `json:"IslandDate"`
}

//모험섬 이름 및 보상
type IslandElement struct {
	Name   string `json:"Name"`
	Reward string `json:"Reward"`
}
*/

//islandDate의 " 시작" 부분과 "\x00"(null) 제거.
//client에서 islandDate를 받자마자 타이머로 환산해 처리.
func ProcessIslandDate(islandDate string) string {
	re, err := regexp.Compile(" 시작")
	if err != nil {
		log.Fatal(err)
	}
	islandDate = re.ReplaceAllString(islandDate, "")

	//"\x00"을 제거
	re, err = regexp.Compile("\x00")
	if err != nil {
		log.Fatal(err)
	}
	islandDate = re.ReplaceAllString(islandDate, "")

	return islandDate
}
