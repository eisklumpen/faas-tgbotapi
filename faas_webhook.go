package tgbotapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (bot *BotAPI) GetHandlerFuncForWebhook(botAccessToken string, onUpdate func(update *Update) error) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		queryBat := req.URL.Query().Get("token")
		if queryBat != botAccessToken {
			fmt.Println("Ignoring Call. Bot Access Token missing/invalid")
			res.Write([]byte("Invalid bot access token"))
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			if r := recover(); r != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Header().Set("Content-Type", "application/json")
				errMsg, _ := json.Marshal(map[string]string{"error": "Paniced. See Server logs."})
				_, _ = res.Write(errMsg)
				log.Println("Recovered from panic in webhook handler func. Error: ", r)
			}
		}()
		update, err := bot.HandleUpdate(req)
		if err != nil {
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			res.WriteHeader(http.StatusBadRequest)
			res.Header().Set("Content-Type", "application/json")
			_, _ = res.Write(errMsg)
			return
		}

		err = onUpdate(update)
		if err != nil {
			log.Printf("Failed to handle update. onUpdate Returned %q", err.Error())
			errMsg, _ := json.Marshal(map[string]string{"error": err.Error()})
			res.WriteHeader(http.StatusInternalServerError)
			res.Header().Set("Content-Type", "application/json")
			_, _ = res.Write(errMsg)
			return
		}
		res.WriteHeader(200)
	}
}
