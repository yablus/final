package email

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/yablus/final/config"
	"github.com/yablus/final/internal/functions"
	"github.com/yablus/final/internal/models"
)

const (
	dataPath = config.DataPath
	showData = config.ShowDataInLogs
)

var fileEmailDataName = config.FileEmailDataName

func makeEmailData() []models.EmailData {
	fmt.Println("================")
	log.Println("Запущен сервис Email")

	var fileEmailData = dataPath + fileEmailDataName
	bufEmailData := functions.GetDataFromFile(fileEmailData)
	if bufEmailData == nil {
		log.Println("Services - Email:", "Getting error: empty data")
		return nil
	}
	log.Printf("Файл %s прочитан\n", fileEmailDataName)

	sep := "\n"
	slice := strings.Split(strings.Trim(string(bufEmailData), sep), sep)
	log.Printf("Данные Email сервиса получены: Всего: %d элементов\n", len(slice))

	data := make([]models.EmailData, 0)
	for _, value := range slice {
		values := strings.Split(value, ";")
		if len(values) != 3 {
			log.Printf("Элемент [%s] удален: Некорректное количество параметров (%d)\n", value, len(values))
			continue
		}
		dt, err := strconv.Atoi(values[2])
		if err != nil || dt == 0 {
			log.Printf("Элемент [%s] удален: Некорректное значение времени доставки (%sms)\n", value, values[2])
			continue
		}
		dataStruct := models.EmailData{
			Country:      values[0],
			Provider:     values[1],
			DeliveryTime: dt,
		}
		if !functions.IsValidEmailData(dataStruct, value) {
			continue
		}
		data = append(data, dataStruct)
	}
	if len(slice) != len(data) {
		log.Printf("Данные Email сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - Email:", err)
			return nil
		}
		log.Println("Исправленные данные Email сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return data
}

func formatEmailData(data []models.EmailData) map[string][][]models.EmailData {
	if data == nil {
		log.Println("Services - Email - formatEmailData:", "Formatting error: empty data")
		return nil
	}
	iso3166data := functions.GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	if iso3166data == nil {
		log.Println("Services - Email - formatEmailData:", "Formatting error: empty iso3166data")
		return nil
	}

	// создадим карту map[string]string key: Iso3166.Country, value: Iso3166.Code
	mapCountries := make(map[string]string)
	for _, v := range iso3166data {
		mapCountries[strings.Trim(v.Country, "\r")] = v.Code
	}
	// создадим карту map[string][]EmailData key: Название страны, value: все элементы сервиса соответствующие стране
	mapDataCountries := make(map[string][]models.EmailData)
	for _, value := range data {
		key := func(m map[string]string, s string) string {
			for k, v := range m {
				if v == s {
					return k
				}
			}
			return ""
		}
		mapDataCountries[key(mapCountries, value.Country)] = append(mapDataCountries[key(mapCountries, value.Country)], value)
	}
	// создадим карту map[string][][]EmailData key: название страны, value: Два отредактированных слайса данных
	mapOut := make(map[string][][]models.EmailData)
	for k, v := range mapDataCountries {
		value := func(d []models.EmailData) [][]models.EmailData {
			r := make([][]models.EmailData, 0)
			sort.SliceStable(d, func(i, j int) bool {
				return d[i].DeliveryTime < d[j].DeliveryTime
			})
			var f = make([]models.EmailData, 3)
			copy(f, d)
			sort.SliceStable(d, func(i, j int) bool {
				return d[i].DeliveryTime > d[j].DeliveryTime
			})
			var s = make([]models.EmailData, 3)
			copy(s, d)
			r = append(r, f, s)
			return r
		}
		mapOut[k] = value(v)
	}
	log.Println("Данные Email сервиса обновлены и отсортированы")

	if showData {
		jsonOut, err := json.Marshal(mapOut)
		if err != nil {
			log.Println("Services - Email - formatEmailData:", err)
			return nil
		}
		log.Println("Корректные данные Email сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return mapOut
}

func GetEmailData() map[string][][]models.EmailData {
	return formatEmailData(makeEmailData())
}
