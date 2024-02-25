package sms

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/usecase/config"
	"github.com/yablus/final/internal/usecase/functions"
)

var (
	showData        = config.B("b_showDataInLogs")
	showFinalData   = config.B("b_showFinalDataInLogs")
	fileSMSDataName = config.S("df_fileSMSDataName")
	dataPath        = config.S("p_dataPath")
)

func makeSMSData() []models.SMSData {
	log.Println("Запущен сервис SMS")

	var fileSMSData = dataPath + fileSMSDataName
	bufSMSData := functions.GetDataFromFile(fileSMSData)
	if bufSMSData == nil {
		log.Println("Services - SMS:", "Getting error: empty data")
		return nil
	}
	log.Printf("Файл %s прочитан\n", fileSMSDataName)

	sep := "\n"
	slice := strings.Split(strings.Trim(string(bufSMSData), sep), sep)
	log.Printf("Данные SMS сервиса получены: Всего: %d элементов\n", len(slice))

	data := make([]models.SMSData, 0)
	for _, value := range slice {
		values := strings.Split(value, ";")
		if len(values) != 4 {
			log.Printf("Элемент [%s] удален: Некорректное количество параметров (%d)\n", value, len(values))
			continue
		}
		dataStruct := models.SMSData{
			Country:      values[0],
			Bandwidth:    values[1],
			ResponseTime: values[2],
			Provider:     values[3],
		}
		if !functions.IsValidSMSData(dataStruct, value) {
			continue
		}
		data = append(data, dataStruct)
	}
	if len(slice) != len(data) {
		log.Printf("Данные SMS сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - SMS:", err)
			return nil
		}
		log.Println("Исправленные данные SMS сервиса (в JSON):")
		fmt.Println(string(jsonOut))
		fmt.Println()
	}

	return data
}

func formatSMSData(data []models.SMSData) [][]models.SMSData {
	if data == nil {
		log.Println("Services - SMS - formatSMSData:", "Formatting error: empty data")
		fmt.Println("----------------------------")
		return nil
	}
	iso3166data := functions.GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	if iso3166data == nil {
		log.Println("Services - SMS - formatSMSData:", "Formatting error: empty iso3166data")
		fmt.Println("----------------------------")
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
	log.Println("Данные SMS сервиса отсортированы и подготовлены для api")

	if showFinalData {
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Println("Services - SMS - formatSMSData:", err)
			fmt.Println("----------------------------")
			return nil
		}
		log.Println("Корректные данные SMS сервиса (в JSON):")
		fmt.Println(string(jsonOut))
	}

	fmt.Println("----------------------------")
	return out
}

func GetSMSData() [][]models.SMSData {
	return formatSMSData(makeSMSData())
}
