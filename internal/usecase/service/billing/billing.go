package billing

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/usecase/config"
	"github.com/yablus/final/internal/usecase/functions"
)

var (
	showFinalData       = config.B("b_showFinalDataInLogs")
	fileBillingDataName = config.S("df_fileBillingDataName")
	dataPath            = config.S("p_dataPath")
)

func makeBillingData() models.BillingData {
	log.Println("Запущен сервис Billing")

	var fileBillingData = dataPath + fileBillingDataName
	bufBillingData := functions.GetDataFromFile(fileBillingData)
	if bufBillingData == nil {
		log.Println("Services - Billing:", "Getting error: empty data")
		fmt.Println("----------------------------")
		return models.BillingData{}
	}
	log.Printf("Файл %s прочитан\n", fileBillingDataName)

	sep := "\n"
	result := strings.Trim(strings.TrimSpace(string(bufBillingData)), sep)

	var correctResult string
	if len(result) != 8 {
		for i := len(result); i < 8; i++ {
			correctResult = "0" + correctResult
		}
		correctResult += result
	}
	var numberMask uint8
	for i := len(result) - 1; i >= 0; i-- {
		temp := uint8(math.Pow(float64((result[i]-48)*2), float64(len(result)-i-1)))
		if i == len(result)-1 && result[i]-48 == 0 {
			temp = 0
		}
		numberMask += temp
	}
	log.Printf("Данные Billing сервиса получены: %s [%s] = %d\n", result, correctResult, numberMask)

	dataBool := make([]bool, 0)
	for i := len(result) - 1; i >= 0; i-- {
		if numberMask&(1<<i) != 0 {
			dataBool = append(dataBool, true)
		} else {
			dataBool = append(dataBool, false)
		}
	}
	data := models.BillingData{
		CreateCustomer: dataBool[0],
		Purchase:       dataBool[1],
		Payout:         dataBool[2],
		Recurring:      dataBool[3],
		FraudControl:   dataBool[4],
		CheckoutPage:   dataBool[5],
	}
	log.Printf("Данные Billing сервиса сформированы для api")

	if showFinalData {
		jsonOut, err := json.Marshal(data)
		if err != nil {
			log.Println("Services - Billing:", err)
			fmt.Println("----------------------------")
			return models.BillingData{}
		}
		log.Println("Данные Billing сервиса (в JSON):")
		fmt.Println(string(jsonOut))
	}

	fmt.Println("----------------------------")
	return data

}

func GetBillingData() models.BillingData {
	return makeBillingData()
}
