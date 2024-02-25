package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	configPath = "files/config.ini"
	mDefault   = getDefaultData()
	mConfig    = getConfigData(configPath)
)

func getDefaultData() map[string]string {
	return map[string]string{
		"b_showDataInLogs":       "false",
		"b_showFinalDataInLogs":  "false",
		"df_fileBillingDataName": "billing.data",
		"df_fileEmailDataName":   "email.data",
		"df_fileSMSDataName":     "sms.data",
		"df_fileVoiceDataName":   "voice.data",
		"du_urlIncidentDataPath": "http://127.0.0.1:8383/accendent",
		"du_urlMMSDataPath":      "http://127.0.0.1:8383/mms",
		"du_urlSupportDataPath":  "http://127.0.0.1:8383/support",
		"p_filesPath":            "files/",
		"p_dataPath":             "simulator/",
	}
}

func getConfigData(path string) map[string]string {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		CreateConfigFile()
	}
	f, err := os.Open(path)
	if err != nil {
		log.Println("Config - getConfigData:", err)
		return nil
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		log.Println("Config - getConfigData:", err)
		return nil
	}
	bufStr := strings.Split(string(buf), "\n")
	m := make(map[string]string)
	for _, v := range bufStr {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "#") {
			continue
		}
		if !strings.Contains(v, "=") {
			continue
		}
		sl := strings.Split(v, "=")
		key := strings.TrimSpace(strings.ReplaceAll(sl[0], " ", ""))
		val := strings.TrimSpace(strings.ReplaceAll(sl[1], " ", ""))
		if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
			val = strings.Trim(val, "\"")
		}
		m[key] = val
	}
	return m
}

func CreateConfigFile() {
	f, err := os.Create(configPath)
	if err != nil {
		log.Println("Config - createConfigFile:", err)
	}
	defer f.Close()
	m := mDefault
	mSortedKeys := make([]string, 0, len(m))
	for k := range m {
		mSortedKeys = append(mSortedKeys, k)
	}
	sort.Strings(mSortedKeys)
	for _, k := range mSortedKeys {
		v := m[k]
		if v != "true" && v != "false" {
			v = fmt.Sprintf("\"%s\"", v)
		}
		s := fmt.Sprintf("%s = %s\n", k, v)
		_, err = f.WriteString(s)
		if err != nil {
			log.Println("Config - createConfigFile:", err)
			continue
		}
	}
}

func S(k string) string {
	if mConfig[k] != "" {
		return mConfig[k]
	}
	return mDefault[k]
}
func B(k string) bool {
	if mConfig[k] == "true" || mConfig[k] == "false" {
		return toBool(mConfig[k])
	}
	return toBool(mDefault[k])
}

func toBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}
}
