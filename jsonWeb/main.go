package main

//
// Uso:  go run main.go quito guayaquil "new york" bogota "la paz" lima santiago riobamba ibarra latacunga boston
// Salida: City BOSTON Country US Lat 42.360253 Lon -71.058291 Temp -12.050000 Feels_like -19.050000
//

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type cityData struct {
	city      string
	country   string
	lat       string
	lon       string
	temp      string
	feelsLike string
}

const APIKEY = ""

//func getContent(bytes []byte, expIni string, expEnd string) string {
//
//	reIni := regexp.MustCompile(expIni + `.*,` + expEnd)
//	reNum := regexp.MustCompile(`(\+|-)?\d+(\.)?\d+`)
//	return string(reNum.Find(reIni.Find(bytes)))
//
//}

func getWebBytes(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		url,
		nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}

func queryCityWeather(city string, client *http.Client) (cityData, error) {

	cd := cityData{}
	cityUrl := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", city, APIKEY)
	cityBytes, err := getWebBytes(client, cityUrl)
	if err != nil {
		return cd, err
	}

	var coordinatesResult []map[string]interface{}
	err = json.Unmarshal(cityBytes, &coordinatesResult)
	if err != nil {
		return cd, err
	}

	country := fmt.Sprintf("%s", coordinatesResult[0]["country"])
	lat := fmt.Sprintf("%f", coordinatesResult[0]["lat"])
	lon := fmt.Sprintf("%f", coordinatesResult[0]["lon"])

	weatherUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, APIKEY)

	weatherBytes, err := getWebBytes(client, weatherUrl)
	if err != nil {
		return cd, err
	}

	var weatherResult map[string]interface{}
	err = json.Unmarshal(weatherBytes, &weatherResult)
	if err != nil {
		return cd, err
	}

	weather := weatherResult["main"].(map[string]interface{})
	temp := fmt.Sprintf("%f", weather["temp"])
	feelsLike := fmt.Sprintf("%f", weather["feels_like"])
	cd.city = city
	cd.country = country
	cd.lat = lat
	cd.lon = lon
	cd.temp = temp
	cd.feelsLike = feelsLike
	return cd, nil
}

func main() {
	cities := []string{}
	ch := make(chan cityData)

	if len(os.Args) < 2 {
		fmt.Println("Ingrese nombres de ciudades...")
		fmt.Println("Uso: go run main.go quito guayaquil \"new york\" bogota \"la paz\" lima santiago riobamba ibarra latacunga boston")
		return
	}

	for _, param := range os.Args[1:] {

		city := strings.ToUpper(param)
		cities = append(cities, city)
	}

	client := &http.Client{Timeout: 30000 * time.Millisecond}

	go func() {
		for _, city := range cities {

			cd, err := queryCityWeather(city, client)
			if err != nil {
				log.Fatal(err)
			}
			ch <- cd
		}
	}()

	for _, city := range cities {
		cd := <-ch
		fmt.Printf("===================== Weather: %s======================\n", city)
		fmt.Println("City", cd.city, "Country", cd.country, "Lat", cd.lat, "Lon", cd.lon, "Temp", cd.temp, "Feels_like", cd.feelsLike)
		fmt.Println("======================================")
	}

	fmt.Println("Fin del programa")
}
