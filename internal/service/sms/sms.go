package sms

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/yablus/final/config"
	"github.com/yablus/final/internal/functions"
	"github.com/yablus/final/internal/models"
)

const (
	dataPath = config.DataPath
	showData = config.ShowDataInLogs
)

var fileSMSDataName = config.FileSMSDataName

func makeSmsData() []models.SMSData {
	fmt.Println("================")
	log.Println("Запущен сервис SMS")

	var fileSMSData = dataPath + fileSMSDataName
	bufSMSData := functions.GetDataFromFile(fileSMSData)
	if bufSMSData == nil {
		log.Fatalln("Services - SMS:", "Getting error: empty data")
	}
	log.Printf("Файл [%s] прочитан\n", fileSMSDataName)

	sep := "\n"
	slice := strings.Split(strings.Trim(string(bufSMSData), sep), sep)
	log.Printf("Данные SMS сервиса получены: Всего: %d строк\n", len(slice))

	data := make([]models.SMSData, 0)
	for _, value := range slice {
		values := strings.Split(value, ";")
		if len(values) != 4 {
			log.Printf("Строка [%s] удалена: Некорректное количество параметров (%d)\n", value, len(values))
			continue
		}
		dataStruct := models.SMSData{
			Country:      values[0],
			Bandwidth:    values[1],
			ResponseTime: values[2],
			Provider:     values[3],
		}
		if !functions.IsValidSMSDataStruct(dataStruct, value) {
			continue
		}
		data = append(data, dataStruct)
	}
	if len(slice) != len(data) {
		log.Printf("Данные SMS сервиса откорректированы: Всего: %d строк\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - SMS:", err)
			return nil
		}
		log.Println("Исправленный дамп данных SMS сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return data
}

func formatSMSData(data []models.SMSData) [][]models.SMSData {
	if data == nil {
		log.Println("Services - SMS - formatSMSData:", "Formatting error: empty data")
		return nil
	}
	iso3166data := functions.GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	if iso3166data == nil {
		log.Println("Services - SMS - formatSMSData:", "Formatting error: empty data")
		return nil
	}
	for i := 0; i < len(data); i++ {
		for _, iso3166 := range iso3166data {
			if data[i].Country == iso3166.Code {
				data[i].Country = strings.Trim(iso3166.Country, "\r")
			}
		}
	}
	var out = make([][]models.SMSData, 0)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Provider < data[j].Provider
	})
	var sortByProvider = make([]models.SMSData, len(data))
	copy(sortByProvider, data)
	out = append(out, sortByProvider)
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Country < data[j].Country
	})
	out = append(out, data)
	log.Println("Данные SMS сервиса обновлены и отсортированы")

	if showData {
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Println("Services - SMS - formatSMSData:", err)
			return nil
		}
		log.Println("Корректный дамп данных SMS сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return out
}

func GetSMSData() [][]models.SMSData {
	return formatSMSData(makeSmsData())
}
