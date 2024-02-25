package service

import (
	"github.com/yablus/final/internal/models"
	"github.com/yablus/final/internal/usecase/service/billing"
	"github.com/yablus/final/internal/usecase/service/email"
	"github.com/yablus/final/internal/usecase/service/incident"
	"github.com/yablus/final/internal/usecase/service/mms"
	"github.com/yablus/final/internal/usecase/service/sms"
	"github.com/yablus/final/internal/usecase/service/support"
	"github.com/yablus/final/internal/usecase/service/voice"
)

type DataService interface {
	GetResultData() models.ResultT
	GetResultSetT() models.ResultSetT
	GetSMS() [][]models.SMSData
	GetMMS() [][]models.MMSData
	GetVoice() []models.VoiceCallData
	GetEmail() map[string][][]models.EmailData
	GetBilling() models.BillingData
	GetSupport() []int
	GetIncident() []models.IncidentData
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
		SMS:      sms.GetSMSData(),
		MMS:      mms.GetMMSData(),
		Voice:    voice.GetVoiceData(),
		Email:    email.GetEmailData(),
		Billing:  billing.GetBillingData(),
		Support:  support.GetSupportData(),
		Incident: incident.GetIncidentData(),
	}
}

func (u *Data) GetResultData() models.ResultT {
	var result models.ResultT
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

func (u *Data) GetVoice() []models.VoiceCallData {
	return u.Voice
}

func (u *Data) GetEmail() map[string][][]models.EmailData {
	return u.Email
}

func (u *Data) GetBilling() models.BillingData {
	return u.Billing
}

func (u *Data) GetSupport() []int {
	return u.Support
}

func (u *Data) GetIncident() []models.IncidentData {
	return u.Incident
}
