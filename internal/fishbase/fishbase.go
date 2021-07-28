package fishbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const baseURL = "https://fishbase.ropensci.org/"

type HeartbeatStatus string

const (
	HeartbeatStatusOperational HeartbeatStatus = "Operational"
	HeartbeatStatusDegraded    HeartbeatStatus = "Degraded"
	HeartbeatStatusDown        HeartbeatStatus = "Down"
)

type Heartbeat struct {
	Status HeartbeatStatus
}

type speciesInfo struct {
	Data []struct {
		SpecCode  int    `json:"SpecCode,omitempty"`
		Genus     string `json:"Genus,omitempty"`
		Species   string `json:"Species,omitempty"`
		Subfamily string `json:"Subfamily,omitempty"`
		Dangerous string `json:"Dangerous,omitempty"`
		Image     string `json:"Image,omitempty"`
	} `json:"data,omitempty"`
}

type ecosystemInfo struct {
	Data []struct {
		SpecCode int    `json:"SpecCode,omitempty"`
		Name     string `json:"EcosystemName,omitempty"`
		Type     string `json:"EcosystemType,omitempty"`
		Location string `json:"Location,omitempty"`
		Salinity string `json:"Salinity,omitempty"`
		Climate  string `json:"Climate,omitempty"`
	} `json:"data,omitempty"`
}

type Ecosystem struct {
	Name     string
	Type     string
	Location string
	Salinity string
	Climate  string
}

type Fish struct {
	SpecCode  int
	Genus     string
	Species   string
	Subfamily string
	Dangerous string
	Image     string
	Ecosystem Ecosystem
}

func GetHeartbeat() (Heartbeat, error) {
	h := Heartbeat{}

	req, err := http.NewRequest("GET", baseURL+"heartbeat", nil)
	if err != nil {
		return h, errors.Wrap(err, "unable to create request")
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return h, errors.Wrap(err, "unable to query fishbase")
	}

	defer rsp.Body.Close()

	h.Status = HeartbeatStatusDown
	if rsp.StatusCode < 400 {
		h.Status = HeartbeatStatusOperational
	}

	return h, nil
}

// GetDetails takes a fishs genus and species, and queries fishbase to collect information about that fish.
func GetDetails(genus string, species string) (Fish, error) {
	fish := Fish{}

	si, err := getSpeciesInformation(genus, species)
	if err != nil {
		return fish, errors.Wrap(err, "unable to get species information")
	}

	ei, err := getEcosystemInformation(strconv.Itoa(si.Data[0].SpecCode))
	if err != nil {
		return fish, errors.Wrap(err, "unable to get ecosystem information")
	}

	fish = Fish{
		SpecCode:  si.Data[0].SpecCode,
		Genus:     si.Data[0].Genus,
		Species:   si.Data[0].Species,
		Subfamily: si.Data[0].Subfamily,
		Dangerous: si.Data[0].Dangerous,
		Image:     si.Data[0].Image,
		Ecosystem: Ecosystem{
			Name:     ei.Data[0].Name,
			Type:     ei.Data[0].Type,
			Location: ei.Data[0].Location,
			Salinity: ei.Data[0].Salinity,
			Climate:  ei.Data[0].Climate,
		},
	}

	return fish, nil
}

func getSpeciesInformation(genus string, species string) (speciesInfo, error) {
	si := speciesInfo{}

	req, err := http.NewRequest("GET", baseURL+"species", nil)
	if err != nil {
		return si, errors.Wrap(err, "unable to create request")
	}

	q := req.URL.Query()
	q.Add("limit", "1")
	q.Add("Genus", genus)
	q.Add("Species", species)

	req.URL.RawQuery = q.Encode()

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return si, errors.Wrap(err, "unable to query fishbase")
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return si, errors.Wrap(err, "unable to read response body")
	}

	if err := json.Unmarshal(body, &si); err != nil {
		return si, errors.Wrap(err, "unable to unmarshal response")
	}

	if len(si.Data) == 0 {
		return si, fmt.Errorf("Genus '%s' and Species '%s' returned no results", genus, species)
	}

	return si, nil
}

func getEcosystemInformation(specCode string) (ecosystemInfo, error) {
	ei := ecosystemInfo{}

	req, err := http.NewRequest("GET", baseURL+"ecosystem", nil)
	if err != nil {
		return ei, errors.Wrap(err, "unable to create request")
	}

	q := req.URL.Query()
	q.Add("limit", "1")
	q.Add("SpecCode", specCode)

	req.URL.RawQuery = q.Encode()

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ei, errors.Wrap(err, "unable to query fishbase")
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return ei, errors.Wrap(err, "unable to read response body")
	}

	if err := json.Unmarshal(body, &ei); err != nil {
		return ei, errors.Wrap(err, "unable to unmarshal response")
	}

	if len(ei.Data) == 0 {
		return ei, fmt.Errorf("SpecCode '%s' returned no results", specCode)
	}

	return ei, nil
}
