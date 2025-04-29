package plug

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Plug struct {
	ID          int
	Name        string
	IsOn        bool
	PowerUsage  float32
	isSelected  bool
	IsProtected bool
	manager     *PlugManager
}

func (p *Plug) Select() {
	if p.manager != nil {
		p.manager.selectPlug(p)
	}
}

func (p *Plug) Toggle(backendUrl string) {
	if p.IsProtected {
		log.Println("Plug is protected, cannot toggle")
		return
	}

	log.Println("Toggling plug", backendUrl+"/plug/"+strconv.Itoa(p.ID)+"/toggle")
	client := &http.Client{
		Timeout: 1 * time.Second, // Set timeout duration here
	}

	url := backendUrl + "/shelly/" + strconv.Itoa(p.ID) + "/toggle"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return
	}
}
