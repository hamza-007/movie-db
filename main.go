package main

import (
	// "encoding/json"
	// "fmt"
	// "io"
	// "net/http"
	// "os"
	"log"

	cli "movies/cli"
)

func main() {

	// var result any

	// url := fmt.Sprintf("https://ott-details.p.rapidapi.com/advancedsearch?sort=latest&page=%d", 3)
	// req, _ := http.NewRequest("GET", url, nil)

	// req.Header.Add("X-RapidAPI-Key", "8b5275a2eamshdaf78f129e9e2c7p107edajsn94134e2d54a0")
	// req.Header.Add("X-RapidAPI-Host", "ott-details.p.rapidapi.com")

	// res, _ := http.DefaultClient.Do(req)

	// defer res.Body.Close()
	// body, _ := io.ReadAll(res.Body)
	// err := json.Unmarshal(body, result)
	// if err != nil {
	// 	return
	// }

	// fmt.Println(result)
	// f, err := os.OpenFile("advanced.json", os.O_RDWR, 0644)
	// if err != nil {
	// 	return
	// }

	// n, err := f.Write(body)
	// fmt.Println(n)
	// if err != nil {
	// 	return
	// }

	// n, err = f.WriteString("\n")

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}

}
