package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type ResponseData struct {
	Image  string `json:"image"`
	Answer string `json:"answer"`
}

func main() {

	apiUrl := "https://yesno.wtf/api"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			r.ParseForm()
			question := r.Form.Get("question")

			if strings.Contains(question, "?") {
				response, err := http.Get(apiUrl)
				if err != nil {
					fmt.Println("GET isteği gönderilemedi:", err)
					return
				}
				defer response.Body.Close()

				responseData, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Println("API yanıtı okunamadı:", err)
					return
				}

				var data ResponseData
				if err := json.Unmarshal(responseData, &data); err != nil {
					fmt.Println("Yanıt çözülemedi:", err)
					return
				}

				t, err := template.ParseFiles("template.html")
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				t.Execute(w, data)

				return
			}
		}

		if r.Method == http.MethodGet {

			t, err := template.ParseFiles("template.html")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			t.Execute(w, nil)
		}
	})

	http.ListenAndServe(":8080", nil)
}
