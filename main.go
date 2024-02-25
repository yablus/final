package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/yablus/final/internal/handlers"
	"github.com/yablus/final/internal/usecase/service"
)

func main() {
	echoStart()
	c := handleConnection()
	echoConnection()
	http.ListenAndServe(":8282", c)
}

func handleConnection() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})
	h := &handlers.ServiceHandler{Data: service.NewService()}
	r.Mount("/api", MountRoute(h))
	//r.Mount("/test", MountRouteTests(h))
	return r
}

func MountRoute(h *handlers.ServiceHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetData) // GET /api
	//r.Get("/", test.HandleTest) // GET /api *для тестирования в веб-сервисе. Настройки в test/test.go

	r.Get("/sms", h.GetSMSData)           // GET /api/sms
	r.Get("/mms", h.GetMMSData)           // GET /api/mms
	r.Get("/voice", h.GetVoiceData)       // GET /api/voice
	r.Get("/email", h.GetEmailData)       // GET /api/email
	r.Get("/billing", h.GetBillingData)   // GET /api/billing
	r.Get("/support", h.GetSupportData)   // GET /api/support
	r.Get("/incident", h.GetIncidentData) // GET /api/incident
	return r
}

func MountRouteTests(h *handlers.ServiceHandler) chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.GetData)                     // GET /test *Для тестирования при подмене handler'а
	r.Get("/ResultSetT", h.GetResultSetTData) // GET /test/ResultSetT
	return r
}

// -------- Оформим чуть-чуть...

var clear map[string]func()

const clearConsole = false

func init() {
	if clearConsole {
		clear = make(map[string]func())
		clear["linux"] = func() {
			cmd := exec.Command("clear") //Linux
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
		clear["windows"] = func() {
			cmd := exec.Command("cmd", "/c", "cls") //Windows
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
}

func echoStart() {
	if clearConsole {
		value, ok := clear[runtime.GOOS]
		if ok {
			value()
		}
	}
	fmt.Println("+--------------------------+")
	fmt.Println("|     Финальная работа     |")
	fmt.Println("|  курса \"Go-разработчик\"  |")
	fmt.Println("|      от Лагунова Е.      |")
	fmt.Println("+--------------------------+")
	fmt.Println()
	time.Sleep(2 * time.Second)
}

func echoConnection() {
	log.Println("Api доступен: http://127.0.0.1:8282/api")
	fmt.Println("----------------------------")
}
