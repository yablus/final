package service

import (
	"encoding/json"
	"log"

	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/service/mms"
	"github.com/yablus/final/internal/service/sms"
	"github.com/yablus/final/test"
)

func marshSMS(b []byte) [][]models.SMSData {
	var data [][]models.SMSData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshSMS:", err)
	}
	return data
}

func marshMMS(b []byte) [][]models.MMSData {
	var data [][]models.MMSData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshMMS:", err)
	}
	return data
}

func marshVC(b []byte) []models.VoiceCallData {
	var data []models.VoiceCallData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshVC:", err)
	}
	return data
}

func marshEmail(b []byte) map[string][][]models.EmailData {
	var data map[string][][]models.EmailData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshEmail:", err)
	}
	return data
}

func marshBil(b []byte) models.BillingData {
	var data models.BillingData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshBil:", err)
	}
	return data
}

func marshS(b []byte) []int {
	var data []int
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshS:", err)
	}
	return data
}

func marshIn(b []byte) []models.IncidentData {
	var data []models.IncidentData
	err := json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalln("Service - marshIn:", err)
	}
	return data
}

type DataService interface {
	GetResultData() models.ResultT
	GetResultSetT() models.ResultSetT
	GetSMS() [][]models.SMSData
	GetMMS() [][]models.MMSData
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
		MMS: mms.GetMMSData(),
		//MMS:      marshMMS(test.TestResponseMMS),
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
	if u.SMS != nil && u.MMS != nil && u.Voice != nil && u.Email != nil && u.Support != nil && u.Incident != nil {
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

func (u *Data) GetMMS() [][]models.MMSData {
	return u.MMS
}
