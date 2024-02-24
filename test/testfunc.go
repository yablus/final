package test

import (
	"encoding/json"
	"log"

	"github.com/yablus/final/internal/service/sms"
)

// Тестирование web-сервиса

func GetDataFromSMS() []byte {
	jsonOut, err := json.Marshal(sms.GetSMSData())
	if err != nil {
		log.Println("Testing - GetDataFromSMS:", err)
		return nil
	}
	return jsonOut
}
