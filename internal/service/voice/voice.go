package voice

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

var fileVoiceDataName = config.FileVoiceDataName

func makeVoiceData() []models.VoiceCallData {
	fmt.Println("================")
	log.Println("Запущен сервис VoiceCall")

	var fileVoiceData = dataPath + fileVoiceDataName
	bufVoiceData := functions.GetDataFromFile(fileVoiceData)
	if bufVoiceData == nil {
		log.Println("Services - VoiceCall:", "Getting error: empty data")
		return nil
	}
	log.Printf("Файл %s прочитан\n", fileVoiceDataName)

	sep := "\n"
	slice := strings.Split(strings.Trim(string(bufVoiceData), sep), sep)
	log.Printf("Данные VoiceCall сервиса получены: Всего: %d элементов\n", len(slice))

	data := make([]models.VoiceCallData, 0)
	for _, value := range slice {
		values := strings.Split(value, ";")
		if len(values) != 8 {
			log.Printf("Элемент [%s] удален: Некорректное количество параметров (%d)\n", value, len(values))
			continue
		}
		cs, err := strconv.ParseFloat(values[4], 32)
		if err != nil {
			log.Printf("Элемент [%s] удален: Некорректное значение показателя стабильности соединения (%s)\n", value, values[4])
			continue
		}
		ttfb, err := strconv.Atoi(values[5])
		if err != nil {
			log.Printf("Элемент [%s] удален: Некорректное значение TTFB (%s)\n", value, values[5])
			continue
		}
		vp, err := strconv.Atoi(values[6])
		if err != nil {
			log.Printf("Элемент [%s] удален: Некорректное значение показателя чистоты связи (%s)\n", value, values[6])
			continue
		}
		moct, err := strconv.Atoi(values[7])
		if err != nil {
			log.Printf("Элемент [%s] удален: Некорректное значение медианы длительности звонка (%s)\n", value, values[7])
			continue
		}
		dataStruct := models.VoiceCallData{
			Country:             values[0],
			Bandwidth:           values[1],
			ResponseTime:        values[2],
			Provider:            values[3],
			ConnectionStability: float32(cs),
			TTFB:                ttfb,
			VoicePurity:         vp,
			MedianOfCallsTime:   moct,
		}
		if !functions.IsValidVoiceData(dataStruct, value) {
			continue
		}
		data = append(data, dataStruct)
	}
	if len(slice) != len(data) {
		log.Printf("Данные VoiceCall сервиса изменены: Всего: %d элементов\n", len(data))
	}

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - VoiceCall:", err)
			return nil
		}
		log.Println("Исправленные данные VoiceCall сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return data
}

func formatVoiceData(data []models.VoiceCallData) []models.VoiceCallData {
	if data == nil {
		log.Println("Services - VoiceCall - formatVoiceData:", "Formatting error: empty data")
		return nil
	}
	iso3166data := functions.GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	if iso3166data == nil {
		log.Println("Services - VoiceCall - formatVoiceData:", "Formatting error: empty iso3166data")
		return nil
	}
	for i := 0; i < len(data); i++ {
		for _, iso3166 := range iso3166data {
			if data[i].Country == iso3166.Code {
				data[i].Country = strings.Trim(iso3166.Country, "\r")
			}
		}
	}
	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Country < data[j].Country
	})
	log.Println("Данные VoiceCall сервиса обновлены и отсортированы")

	if showData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - VoiceCall - formatVoiceData:", err)
			return nil
		}
		log.Println("Корректные данные VoiceCall сервиса (в JSON):")
		fmt.Println("----------------")
		fmt.Println(string(jsonOut))
		fmt.Println("----------------")
	}

	return data
}

func GetVoiceData() []models.VoiceCallData {
	return formatVoiceData(makeVoiceData())
}
