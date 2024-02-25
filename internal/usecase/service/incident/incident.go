package incident

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"

	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/usecase/config"
	"github.com/yablus/final/internal/usecase/functions"
)

var (
	showData        = config.B("b_showDataInLogs")
	showFinalData   = config.B("b_showFinalDataInLogs")
	urlIncidentData = config.S("du_urlIncidentDataPath")
)

func makeIncidentData() []models.IncidentData {
	log.Println("Запущен сервис Incident")

	bufIncidentData, status := functions.GetDataFromURL(urlIncidentData)
	if bufIncidentData == nil {
		log.Printf("Services - Incident: Getting error: empty data. Сonnection status code: %d\n", status)
		return nil
	}
	log.Printf("Соединениие по адресу %s установлено\n", urlIncidentData)

	var data []models.IncidentData
	err := json.Unmarshal(bufIncidentData, &data)
	if err != nil {
		log.Println("Services - Incident:", err)
		return nil
	}
	log.Printf("Данные Incident сервиса получены: Всего: %d элементов\n", len(data))

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
		log.Printf("Данные Incident сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - Incident:", err)
			return nil
		}
		log.Println("Исправленные данные Incident сервиса (в JSON):")
		fmt.Println(string(jsonOut))
		fmt.Println()
	}

	return data
}

func formatMIncidentData(data []models.IncidentData) []models.IncidentData {
	if data == nil {
		log.Println("Services - Incident - formatIncidentData:", "Formatting error: empty data")
		fmt.Println("----------------------------")
		return nil
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Status < data[j].Status
	})
	log.Println("Данные Incident сервиса отсортированы и подготовлены для api")

	if showFinalData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - Incident - formatIncidentData:", err)
			fmt.Println("----------------------------")
			return nil
		}
		log.Println("Корректные данные Incident сервиса (в JSON):")
		fmt.Println(string(jsonOut))
	}

	fmt.Println("----------------------------")
	return data
}

func GetIncidentData() []models.IncidentData {
	return formatMIncidentData(makeIncidentData())
}
