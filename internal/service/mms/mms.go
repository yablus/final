package mms

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/yablus/final/internal/config"
	"github.com/yablus/final/internal/functions"
	"github.com/yablus/final/internal/models"
)

var (
	showData      = config.B("b_showDataInLogs")
	showFinalData = config.B("b_showFinalDataInLogs")
	urlMMSData    = config.S("du_urlMMSDataPath")
)

func makeMMSData() []models.MMSData {
	log.Println("Запущен сервис MMS")

	bufMMSData, status := functions.GetDataFromURL(urlMMSData)
	if bufMMSData == nil {
		log.Printf("Services - MMS: Getting error: empty data. Сonnection status code: %d\n", status)
		return nil
	}
	log.Printf("Соединениие по адресу %s установлено\n", urlMMSData)

	var data []models.MMSData
	err := json.Unmarshal(bufMMSData, &data)
	if err != nil {
		log.Println("Services - MMS:", err)
		return nil
	}
	log.Printf("Данные MMS сервиса получены: Всего: %d элементов\n", len(data))

	dataCorrected := false
	for i, v := range data {
		count := reflect.ValueOf(v).NumField()
		if count != 4 {
			dataCorrected = true
			data = append(data[:i], data[i+1:]...)
			log.Printf("Элемент %v удален: Некорректное количество параметров (%d)\n", v, count)
			continue
		}
		if !functions.IsValidMMSData(v) {
			dataCorrected = true
			data = append(data[:i], data[i+1:]...)
			continue
		}
	}
	if dataCorrected {
		log.Printf("Данные MMS сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - MMS:", err)
			return nil
		}
		log.Println("Исправленные данные MMS сервиса (в JSON):")
		fmt.Println(string(jsonOut))
		fmt.Println()
	}

	return data
}

func formatMMSData(data []models.MMSData) [][]models.MMSData {
	if data == nil {
		log.Println("Services - MMS - formatMMSData:", "Formatting error: empty data")
		return nil
	}
	iso3166data := functions.GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	if iso3166data == nil {
		log.Println("Services - MMS - formatMMSData:", "Formatting error: empty iso3166data")
		return nil
	}
	for i := 0; i < len(data); i++ {
		for _, iso3166 := range iso3166data {
			if data[i].Country == iso3166.Code {
				data[i].Country = strings.Trim(iso3166.Country, "\r")
			}
		}
	}
	var out = make([][]models.MMSData, 0)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Provider < data[j].Provider
	})
	var sortByProvider = make([]models.MMSData, len(data))
	copy(sortByProvider, data)
	out = append(out, sortByProvider)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Country < data[j].Country
	})
	out = append(out, data)
	log.Println("Данные MMS сервиса отсортированы и подготовлены для api")

	if showFinalData {
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Println("Services - MMS - formatMMSData:", err)
			return nil
		}
		log.Println("Корректные данные MMS сервиса (в JSON):")
		fmt.Println(string(jsonOut))
	}

	fmt.Println("----------------------------")
	return out
}

func GetMMSData() [][]models.MMSData {
	return formatMMSData(makeMMSData())
}
