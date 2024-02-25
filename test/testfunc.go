package test

import (
	"encoding/json"
	"log"

	"github.com/yablus/final/internal/usecase/service/billing"
	"github.com/yablus/final/internal/usecase/service/email"
	"github.com/yablus/final/internal/usecase/service/mms"
	"github.com/yablus/final/internal/usecase/service/sms"
	"github.com/yablus/final/internal/usecase/service/support"
	"github.com/yablus/final/internal/usecase/service/voice"
	//"github.com/yablus/final/internal/usecase/service/incident"
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

func GetDataFromMMS() []byte {
	jsonOut, err := json.Marshal(mms.GetMMSData())
	if err != nil {
		log.Println("Testing - GetDataFromMMS:", err)
		return nil
	}
	return jsonOut
}

func GetDataFromVoice() []byte {
	jsonOut, err := json.Marshal(voice.GetVoiceData())
	if err != nil {
		log.Println("Testing - GetDataFromVoice:", err)
		return nil
	}
	return jsonOut
}

func GetDataFromEmail() []byte {
	jsonOut, err := json.Marshal(email.GetEmailData())
	if err != nil {
		log.Println("Testing - GetDataFromEmail:", err)
		return nil
	}
	return jsonOut
}

func GetDataFromBilling() []byte {
	jsonOut, err := json.Marshal(billing.GetBillingData())
	if err != nil {
		log.Println("Testing - GetDataFromBilling:", err)
		return nil
	}
	return jsonOut
}

func GetDataFromSupport() []byte {
	jsonOut, err := json.Marshal(support.GetSupportData())
	if err != nil {
		log.Println("Testing - GetDataFromSupport:", err)
		return nil
	}
	return jsonOut
}

/*
func GetDataFromIncident() []byte {
	jsonOut, err := json.Marshal(incident.GetIncidentData())
	if err != nil {
		log.Println("Testing - GetDataFromIncident:", err)
		return nil
	}
	return jsonOut
}
*/
