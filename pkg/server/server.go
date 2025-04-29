package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/maxlehmann01/hmon-terminal/pkg/plug"
	"github.com/maxlehmann01/hmon-terminal/pkg/ui"
)

type JSONPlug struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	IsOn        bool    `json:"is_on"`
	IsProtected bool    `json:"is_protected"`
	PowerUsage  float32 `json:"power_usage"`
}

func Start(pm *plug.PlugManager, port int) {
	http.HandleFunc("/plugs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var jsonPlugs []JSONPlug
		if err := json.Unmarshal(body, &jsonPlugs); err != nil {
			http.Error(w, "invalid json format", http.StatusBadRequest)
			return
		}

		selectedPlug := pm.GetSelected()
		selectedID := -1

		if selectedPlug != nil {
			selectedID = selectedPlug.ID
		}

		pm.Clear()
		var foundSelected bool
		var minID int
		first := true

		for _, jsonPlug := range jsonPlugs {
			pm.AddPlug(&plug.Plug{
				ID:          jsonPlug.ID,
				Name:        jsonPlug.Name,
				IsOn:        jsonPlug.IsOn,
				IsProtected: jsonPlug.IsProtected,
				PowerUsage:  jsonPlug.PowerUsage,
			})

			if jsonPlug.ID == selectedID {
				foundSelected = true
			}

			if first || jsonPlug.ID < minID {
				minID = jsonPlug.ID
				first = false
			}
		}

		if foundSelected {
			pm.SelectPlugByID(selectedID)
		} else if !first {
			pm.SelectPlugByID(minID)
		}

		ui.OutputSelectedPlug(pm)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("HTTP server running on :" + strconv.Itoa(port))
	go http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
