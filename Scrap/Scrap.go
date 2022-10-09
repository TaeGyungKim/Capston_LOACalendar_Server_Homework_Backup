/* 스크랩 패키지

크롤링
모험섬 데이터

아이템 데이터

*/

package Item

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//아이템 이름, 고유번호, 가격

//아이템에 따른 가격 및 매물대
type ItemPrice struct {
	Name       string       `json:"Name"`
	Pricechart []Pricechart `json:"Pricechart"`
	Result     string       `json:"Result"`
}

//각 가격
type Pricechart struct {
	Amount string `json:"Amount"`
	Price  string `json:"Price"`
}

type Item struct {
	name string `json:"name"`
	num  string `json:"num"`
}

var ItemMap map[int]Item
var ItemNumber int

/*
func Scrap2() {
	var island Island

	islandData := "http://152.70.248.4:5000/adventureisland/"

	err := json.NewDecoder(res.Body).Decode(&island)
}
*/

func UpdateItemNumber() {

}

//데이터를 외부에서 입력받음
//개선안 : 트리를 이용하여 이름을 추적할 수 있도록하기 - 보류
//
func InputData(data string) (string, error) {
	var item Item
	var err error
	//convertToInt, err := strconv.Atoi(data)

	item.name = data

	item.num, err = item.matching()
	if err != nil {
		log.Fatal(err)

	}

	return item.num, err
}

//아이템 맵에서 name맵 순회하여 num을 찾음.
func (i Item) matching() (string, error) {
	var item Item
	item.name = i.name

	for _, v := range ItemMap {
		if v.name == item.name {
			item.num = v.num
			log.Print("아이템명과 매칭완료, 아이템 번호:", v.num)
			return v.num, nil
		}
	}

	return "0", fmt.Errorf("item not matched")
}

func GetPriceData(ITEM_NUM string) (ItemPrice, error) {
	var itemPrice ItemPrice

	//주소에 고유번호 더해서 검색

	//예제 : 오레하 유물, 6885708
	//6885708
	priceData := "http://152.70.248.4:5000/trade" + ITEM_NUM

	res, err := http.Get(priceData)
	if err != nil {
		log.Println("Get 에러")
		log.Fatal(err)
	}
	defer res.Body.Close()
	r, err := http.NewRequest("GET", "/", res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//*request.body를 디코딩해서 island에 저장
	err = json.NewDecoder(r.Body).Decode(&itemPrice)
	if err != nil {
		log.Fatal(err)
	}

	//출력 로그.
	log.Println(itemPrice)

	return itemPrice, err

}

//미사용
//
//이 링크 참조하여 변환할것.
//https://jeonghwan-kim.github.io/dev/2019/01/18/go-encoding-json.html

//사용할지 모르니 대기.
func tempt() {
	var info string

	priceData := "http://152.70.248.4:5000/trade"

	res, err := http.Get(priceData)
	if err != nil {
		log.Println("Get 에러")
		log.Fatal(err)
	}

	//만일 응답 상태 코드가 200번이 아니라면 오류
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	log.Println("status code:", res.StatusCode)
	log.Println("Body:", res.Body)
	//http의 Body부분을 크롤링해서 data를 얻음
	data, err := goquery.NewDocumentFromReader(res.Body)

	log.Println("data:", data)
	if err != nil {
		log.Fatal(err)
	}

	//data의 필요한 부분만 검색
	//*Document 변수를 string으로 하여 사용하기 위함
	data.Find("").Each(func(_ int, s *goquery.Selection) {
		info = s.Find("").Text()
	})

	log.Println("info:", info)

	//string을 []byte로 변환
	//var convertByte = []byte(info)
	//	island, err := UnmarshalIsland(convertByte)

	//	convertByte, err = island.Marshal()
}
