package support

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/yablus/final/internal/config"
	"github.com/yablus/final/internal/functions"
	"github.com/yablus/final/internal/models"
)

var (
	showData       = config.B("b_showDataInLogs")
	showFinalData  = config.B("b_showFinalDataInLogs")
	urlSupportData = config.S("du_urlSupportDataPath")
)

func makeSupportData() []models.SupportData {
	log.Println("Запущен сервис Support")

	bufSupportData, status := functions.GetDataFromURL(urlSupportData)
	if bufSupportData == nil {
		log.Printf("Services - Support: Getting error: empty data. Сonnection status code: %d\n", status)
		return nil
	}
	log.Printf("Соединениие по адресу %s установлено\n", urlSupportData)

	var data []models.SupportData
	err := json.Unmarshal(bufSupportData, &data)
	if err != nil {
		log.Println("Services - Support:", err)
		return nil
	}
	log.Printf("Данные Support сервиса получены: Всего: %d элементов\n", len(data))

	dataCorrected := false
	for i, v := range data {
		count := reflect.ValueOf(v).NumField()
		if count != 2 {
			dataCorrected = true
			data = append(data[:i], data[i+1:]...)
			log.Printf("Элемент %v удален: Некорректное количество параметров (%d)\n", v, count)
			continue
		}
	}
	if dataCorrected {
		log.Printf("Данные Support сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - Support:", err)
			return nil
		}
		log.Println("Данные Support сервиса (в JSON):")
		fmt.Println(string(jsonOut))
		fmt.Println()
	}
	return data
}

func formatSupportData(data []models.SupportData) []int {
	if data == nil {
		log.Println("Services - Support - formatSupportData:", "Formatting error: empty data")
		fmt.Println("----------------------------")
		return nil
	}
	var count int
	for _, v := range data {
		count += v.ActiveTickets
	}
	var loadIndex int
	switch {
	case count > 16:
		loadIndex = 3
	case count > 8:
		loadIndex = 2
	default:
		loadIndex = 1
	}
	var minlt float64 = 60 / 18
	requestTime := int(float64(count) * minlt)
	log.Println("Данные Support сервиса для api подсчитаны")
	out := []int{loadIndex, requestTime}

	if showFinalData {
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Println("Services - Support - formatSupportData:", err)
			fmt.Println("----------------------------")
			return nil
		}
		log.Println("Данные Support сервиса:")
		fmt.Println(string(jsonOut))
	}

	fmt.Println("----------------------------")
	return out
}

func GetSupportData() []int {
	return formatSupportData(makeSupportData())
}
