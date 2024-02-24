package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yablus/final/internal/service"
)

type ServiceHandler struct {
	Data service.DataService
}

func (h *ServiceHandler) GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	err := json.NewEncoder(w).Encode(h.Data.GetResultData())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Println("Данные переданы:")
}

// Для тестирования:

func (h *ServiceHandler) GetResultSetTData(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(h.Data.GetResultSetT())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Println("Данные ResultSetT переданы:")
}

func (h *ServiceHandler) GetSMSData(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(h.Data.GetSMS())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Println("Данные [][]SMSData переданы:")
}

func (h *ServiceHandler) GetMMSData(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(h.Data.GetMMS())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Println("Данные [][]MMSData переданы:")
}

func (h *ServiceHandler) GetVoiceData(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(h.Data.GetVoice())
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println("Internal error")
		return
	}
	log.Println("Данные []VoiceCallData переданы:")
}
