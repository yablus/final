package service

import (
	"encoding/json"
	"log"

	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/service/sms"
	"github.com/yablus/final/test"
)

func marshSMS(b []byte) [][]models.SMSData {
	var data [][]models.SMSData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshSMS:", err)
	}
	return data
}

func marshMMS(b []byte) [][]models.MMSData {
	var data [][]models.MMSData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshMMS:", err)
	}
	return data
}

func marshVC(b []byte) []models.VoiceCallData {
	var data []models.VoiceCallData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshVC:", err)
	}
	return data
}

func marshEmail(b []byte) map[string][][]models.EmailData {
	var data map[string][][]models.EmailData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshEmail:", err)
	}
	return data
}

func marshBil(b []byte) models.BillingData {
	var data models.BillingData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshBil:", err)
	}
	return data
}

func marshS(b []byte) []int {
	var data []int
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshS:", err)
	}
	return data
}

func marshIn(b []byte) []models.IncidentData {
	var data []models.IncidentData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("marshIn:", err)
	}
	return data
}

type DataService interface {
	GetResultData() models.ResultT
	GetResultSetT() models.ResultSetT
	GetSMS() [][]models.SMSData
}

type Data struct {
	SMS      [][]models.SMSData
	MMS      [][]models.MMSData
	Voice    []models.VoiceCallData
	Email    map[string][][]models.EmailData
	Billing  models.BillingData
	Support  []int
	Incident []models.IncidentData
}

func NewService() *Data {
	return &Data{
		SMS: sms.GetSMSData(),
		//SMS:      marshSMS(test.TestResponseSMS),
		MMS:      marshMMS(test.TestResponseMMS),
		Voice:    marshVC(test.TestResponseVoiceCall),
		Email:    marshEmail(test.TestResponseEmail),
		Billing:  marshBil(test.TestResponseBilling),
		Support:  marshS(test.TestResponseSupport),
		Incident: marshIn(test.TestResponseIncident),
	}
}

func (u *Data) GetResultData() models.ResultT {
	var result models.ResultT
	//Проверить все поля
	if u.SMS != nil {
		result.Status = true
		result.Data = u.GetResultSetT()
	} else {
		result.Status = false
		result.Error = "Error on collect data"
	}
	return result
}

func (u *Data) GetResultSetT() models.ResultSetT {
	var result = models.ResultSetT{
		SMS:       u.SMS,
		MMS:       u.MMS,
		VoiceCall: u.Voice,
		Email:     u.Email,
		Billing:   u.Billing,
		Support:   u.Support,
		Incidents: u.Incident,
	}
	return result
}

func (u *Data) GetSMS() [][]models.SMSData {
	return u.SMS
}
