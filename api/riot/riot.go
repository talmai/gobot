package riot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	RIOT_VERSION   = 2.0
	GOBOT_URL_PATH = "http://127.0.0.1:8080/api/robots/riotBot/devices/riot/commands/"
)

// Parses CommandPath returns the result for the given name.
// It returns an error if the path could not be found or
// could not be loaded.
func CommandPath(name string, res http.ResponseWriter, req *http.Request) ([]byte, error) {

	if strings.Contains(name, "riot/sensors") {
		// a.Get("/riot/sensors", a.riot)
		// a.Post("/riot/sensors", a.riot)
		// a.Get("/riot/sensors/:id", a.riot)
		switch len(strings.Split(name, "/")) {
		case 3:
			return sensor(res, req)
		default:
			return sensors(res, req)
		}
	} else {
		if strings.Contains(name, "riot/digital/") {
			if req.Method == "GET" {
				// a.Get("/riot/digital/input/:channel", a.riot)   // channels [0-3]
				// a.Get("/riot/digital/output/:channel", a.riot)  // channels [0-1]
				// a.Get("/riot/digital/relay/:channel", a.riot)   // channels [0-1]
				return readDigitalInput(res, req)
			} else {
				if strings.Contains(name, "riot/digital/output/") {
					return setDigitalOuput(res, req)
				} else {
					// a.Post("/riot/digital/relay/:channel", a.riot)  // channels [0-1]
					return setRelayOuput(res, req)
				}
			}
		} else {
			if strings.Contains(name, "riot/analog") {
				if strings.Contains(name, "output") {
					// a.Get("/riot/analog/output/:channel", a.riot)   // channels [0-1]
					// a.Post("/riot/analog/output/:channel", a.riot)  // channels [0-1]
					return analogOutput(res, req)
				} else {
					// a.Get("/riot/analog/input/:channel", a.riot)    // channels [0-3]
					return readAnalogInput(res, req)
				}
			} else {
				// a.Get("/riot", a.riot)
				return []byte("RIOT version: " + strconv.FormatFloat(RIOT_VERSION, 'f', 6, 64)), nil
			}
		}
	}

	return nil, fmt.Errorf("COMMAND PATH %s INCOMPLETE.", name)
}

func readAnalogInput(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	channel := req.URL.Query().Get(":channel")

	urlCall := GOBOT_URL_PATH + "ReadADCChannel"

	switch channel {
	case "0":
		urlCall += "Zero"
	case "1":
		urlCall += "One"
	case "2":
		urlCall += "Two"
	default:
		urlCall += "Three"
	}

	response, err := http.Get(urlCall)
	if err == nil {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)
		return []byte(buf), nil
	}

	return nil, fmt.Errorf("ERROR OCCURRED %s.", err)
}

func readDigitalInput(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	// channel := req.URL.Query().Get(":channel") // does not matter for now

	response, err := http.Get(GOBOT_URL_PATH + "ReadDigitalInput")

	if err == nil {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)
		return []byte(buf), nil
	}

	return nil, fmt.Errorf("ERROR OCCURRED %s.", err)
}

func setRelayOuput(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	channel := req.URL.Query().Get(":channel")

	req.ParseForm()            //Parse url parameters passed, then parse the response packet for the POST body (request body)
	value := req.Form["value"] // 0 => reset, 1 => set

	urlCall := GOBOT_URL_PATH

	switch value[0] {
	case "0":
		{
			if channel == "0" {
				urlCall += "ResetRelayOutputChannelZero"
			} else {
				urlCall += "ResetRelayOutputChannelOne"
			}
		}
	default:
		{
			if channel == "0" {
				urlCall += "SetRelayOutputChannelZero"
			} else {
				urlCall += "SetRelayOutputChannelOne"
			}
		}
	}

	response, err := http.Get(urlCall)
	if err == nil {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)
		return []byte(buf), nil
	}

	return nil, fmt.Errorf("ERROR OCCURRED %s.", err)
}

func setDigitalOuput(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	channel := req.URL.Query().Get(":channel")

	req.ParseForm()            //Parse url parameters passed, then parse the response packet for the POST body (request body)
	value := req.Form["value"] // 0 => reset, 1 => set

	urlCall := GOBOT_URL_PATH

	switch value[0] {
	case "0":
		{
			if channel == "0" {
				urlCall += "ResetDigitalOutputChannelZero"
			} else {
				urlCall += "ResetDigitalOutputChannelOne"
			}
		}
	default:
		{
			if channel == "0" {
				urlCall += "SetDigitalOutputChannelZero"
			} else {
				urlCall += "SetDigitalOutputChannelOne"
			}
		}
	}

	response, err := http.Get(urlCall)
	if err == nil {
		defer response.Body.Close()
		buf, _ := ioutil.ReadAll(response.Body)
		return []byte(buf), nil
	}

	return nil, fmt.Errorf("ERROR OCCURRED %s.", err)
}

func sensors(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	return nil, fmt.Errorf("SENSOR API NOT YET IMPLEMENTED.")
}

func sensor(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	id := req.URL.Query().Get(":id")
	return nil, fmt.Errorf("THERE IS NO SENSOR REGISTERED UNDER ID %s. SENSOR API NOT YET IMPLEMENTED.", id)
}

func analogOutput(res http.ResponseWriter, req *http.Request) ([]byte, error) {
	return nil, fmt.Errorf("READ/WRITE ANALOG OUTPUT NOT YET IMPLEMENTED.")
}
