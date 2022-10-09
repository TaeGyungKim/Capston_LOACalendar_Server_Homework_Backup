/*
제출 : 2022.06.01.
백업 : 2022.10.09.

캡스톤 프로젝트
컴퓨터 실행하면 당일의 로아 데이터를 불러와서 보여주는 시작프로그램
*/

package main

import (
	Scrap "Capston/Scrap" //Scrap 기능
	Processing "Capston/processing"

	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

//island 관리할 맵 데이터
var IslandMap map[bool]Island

// 아이템 관리할 맵 데이터
var ItemMap map[int]Item

// 맵에 데이터 있는지 확인.
var mapIslandCheck bool

// 아이템 맵 검색에 따른 정보 저장
var itemSearchNum int = 0

//메인 핸들러
func MainHandler() http.Handler {
	IslandMap = make(map[bool]Island)
	ItemMap = make(map[int]Item)

	r := mux.NewRouter()
	r.Handle("/", r)
	//Restful api. (cors)
	r.HandleFunc("/island", GetIslandHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/item", GetItemHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/item", PostItemHandler).Methods(http.MethodPost)
	r.HandleFunc("/item/get", UpdateItemHandler).Methods("POST")
	r.HandleFunc("/item/{num:[0-9]+}", RemoveItemHandler).Methods("DELETE")

	r.Use(mux.CORSMethodMiddleware(r))
	return r
}

//해당 도메인에 데이터를 요청(GET)하면
//island 데이터를 외부에서 스크랩해서 가져옴
// 맵에 해당 데이터를 넣고
// JSON 데이터를 출력
func GetIslandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")

	var island Island
	var err error
	island, err = island.GetData()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mapIslandCheck = true

	IslandMap[mapIslandCheck] = island

	rd.JSON(w, http.StatusOK, island)
}

type Item struct {
	name string `json:"name"`
	num  string `json:"num"`
}

//json형태로 프론트에서 데이터 받아와야함.
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Max-Age", "86400")

}

//
func PostItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("클라이언트로부터 요청을 받아옵니다.")

	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("입력받은 값:", item)

	item.num, err = Scrap.InputData(item.name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound) //404
		return
	}

	itemSearchNum++

	ItemMap[itemSearchNum] = item

	rd.JSON(w, http.StatusOK, item)
}

type Success struct {
	Success bool `json: "success"`
}

// 임시
//db에 데이터 추가하는 기능
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {

}

//임시
//만일 잘못 입력했으면 삭제하는 기능
func RemoveItemHandler(w http.ResponseWriter, r *http.Request) {

}

//메인함수
func main() {
	//서버 띄우기.
	rd = render.New()
	m := MainHandler()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App!!")

	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, n)
	//err := http.ListenAndServe(":3000", n) //test code
	if err != nil {
		panic(err)
	}

}

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

//모험섬 데이터
func (island Island) GetData() (Island, error) {

	//해당 데이터의 json 주소
	islandData := "http://152.70.248.4:5000/adventureisland/"

	//http 획득 시도
	//*response
	res, err := http.Get(islandData)
	if err != nil {
		log.Println("Get 에러")
		log.Fatal(err)
	}
	defer res.Body.Close()

	//*response.body를 이용하여 *request
	r, err := http.NewRequest("GET", "/", res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//*request.body를 디코딩해서 island에 저장
	err = json.NewDecoder(r.Body).Decode(&island)
	if err != nil {
		log.Fatal(err)
	}

	island.IslandDate = Processing.ProcessIslandDate(island.IslandDate)

	//출력 로그.
	log.Println(island)

	return island, err
}
