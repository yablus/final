package functions

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/yablus/final/internal/config"
	"github.com/yablus/final/internal/models"
)

var filesPath = config.S("p_filesPath")

type Iso3166 struct {
	Code    string
	Country string
}

func GetDataFromURL(url string) ([]byte, int) {
	r, err := http.Get(url)
	if err != nil {
		log.Println("Functions - GetDataFromURL:", err)
		return nil, 404
	}
	defer r.Body.Close()
	if r.StatusCode == 200 { // OK
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Functions - GetDataFromURL:", err)
			return nil, 404
		}
		return buf, 200
	}
	return nil, r.StatusCode
}

func GetDataFromFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Functions - GetDataFromFile:", err)
		return nil
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		log.Println("Functions - GetDataFromFile:", err)
		return nil
	}
	return buf
}

func GetAllCountriesFromFile(fileName string) []Iso3166 {
	var filePath = filesPath + fileName
	buf := GetDataFromFile(filePath)
	if buf == nil {
		log.Fatalln("functions - GetAllCountriesFromFile:", "Getting error: empty data")
	}
	dataLong := strings.Split(strings.Trim(strings.TrimSpace(string(buf)), ";"), "\n")
	var data []Iso3166
	for _, v := range dataLong {
		dataShort := strings.Split(v, ";")
		dataIso3166 := Iso3166{
			Code:    dataShort[0],
			Country: dataShort[1],
		}
		data = append(data, dataIso3166)
	}
	return data
}

func GetAllCountryCodes() []string {
	iso3166data := GetAllCountriesFromFile("iso3166-1_alpha-2.data")
	var data []string
	for _, v := range iso3166data {
		data = append(data, v.Code)
	}
	return data
}

func GetAllProvidersFromFile(fileName string) []string {
	var filePath = filesPath + fileName
	buf := GetDataFromFile(filePath)
	if buf == nil {
		log.Fatalln("functions - GetAllProvidersFromFile:", "Getting error: empty data")
	}
	sep := ";"
	data := strings.Split(strings.Trim(strings.TrimSpace(string(buf)), sep), sep)
	return data
}

func SliceContainsString(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

//----------------------

func IsValidSMSData(data models.SMSData, str string) bool {
	if !SliceContainsString(GetAllCountryCodes(), data.Country) {
		log.Printf("Элемент [%s] удален: Код страны %s отсутствует в базе iso3166-1\n", str, data.Country)
		return false
	}
	bandwidthInt, err := strconv.Atoi(data.Bandwidth)
	if err != nil || bandwidthInt < 0 || bandwidthInt > 100 {
		log.Printf("Элемент [%s] удален: Некорректное значение пропускной способности канала (%s)\n", str, data.Bandwidth)
		return false
	}
	_, err = strconv.Atoi(data.ResponseTime)
	if err != nil {
		log.Printf("Элемент [%s] удален: Некорректное значение времени ответа (%sms)\n", str, data.ResponseTime)
		return false
	}
	if !SliceContainsString(GetAllProvidersFromFile("mobile_providers.data"), data.Provider) {
		log.Printf("Элемент [%s] удален: Провайдер %s отсутствует в базе провайдеров\n", str, data.Provider)
		return false
	}
	return true
}

func IsValidMMSData(data models.MMSData) bool {
	if !SliceContainsString(GetAllCountryCodes(), data.Country) {
		log.Printf("Элемент %s удален: Код страны %s отсутствует в базе iso3166-1\n", data, data.Country)
		return false
	}
	bandwidthInt, err := strconv.Atoi(data.Bandwidth)
	if err != nil || bandwidthInt < 0 || bandwidthInt > 100 {
		log.Printf("Элемент %s удален: Некорректное значение пропускной способности канала (%s)\n", data, data.Bandwidth)
		return false
	}
	_, err = strconv.Atoi(data.ResponseTime)
	if err != nil {
		log.Printf("Элемент %s удален: Некорректное значение времени ответа (%sms)\n", data, data.ResponseTime)
		return false
	}
	if !SliceContainsString(GetAllProvidersFromFile("mobile_providers.data"), data.Provider) {
		log.Printf("Элемент %s удален: Провайдер %s отсутствует в базе провайдеров\n", data, data.Provider)
		return false
	}
	return true
}

func IsValidVoiceData(data models.VoiceCallData, str string) bool {
	if !SliceContainsString(GetAllCountryCodes(), data.Country) {
		log.Printf("Элемент [%s] удален: Код страны %s отсутствует в базе iso3166-1\n", str, data.Country)
		return false
	}
	bandwidthInt, err := strconv.Atoi(data.Bandwidth)
	if err != nil || bandwidthInt < 0 || bandwidthInt > 100 {
		log.Printf("Элемент [%s] удален: Некорректное значение пропускной способности канала (%s)\n", str, data.Bandwidth)
		return false
	}
	_, err = strconv.Atoi(data.ResponseTime)
	if err != nil {
		log.Printf("Элемент [%s] удален: Некорректное значение времени ответа (%sms)\n", str, data.ResponseTime)
		return false
	}
	if !SliceContainsString(GetAllProvidersFromFile("voice_providers.data"), data.Provider) {
		log.Printf("Элемент [%s] удален: Провайдер %s отсутствует в базе провайдеров\n", str, data.Provider)
		return false
	}
	return true
}

func IsValidEmailData(data models.EmailData, str string) bool {
	if !SliceContainsString(GetAllCountryCodes(), data.Country) {
		log.Printf("Элемент [%s] удален: Код страны %s отсутствует в базе iso3166-1\n", str, data.Country)
		return false
	}
	if !SliceContainsString(GetAllProvidersFromFile("email_providers.data"), data.Provider) {
		log.Printf("Элемент [%s] удален: Провайдер %s отсутствует в базе провайдеров\n", str, data.Provider)
		return false
	}
	return true
}
