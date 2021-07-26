package fishbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const baseURL = "https://fishbase.ropensci.org/"

type SpeciesInfo struct {
	Data []struct {
		SpecCode  int32  `json:"SpecCode,omitempty"`
		Genus     string `json:"Genus,omitempty"`
		Species   string `json:"Species,omitempty"`
		Subfamily string `json:"Subfamily,omitempty"`
		Dangerous string `json:"Dangerous,omitempty"`
		Image     string `json:"Image,omitempty"`
	} `json:"data,omitempty"`
}

type Fish struct{}

func GetByGenusSpecies(genus string, species string) (*Fish, error) {
	si, err := GetSpeciesInformation(genus, species)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get species information")
	}

	fmt.Printf("%+v \n", si)

	return nil, nil
}

func GetSpeciesInformation(genus string, species string) (SpeciesInfo, error) {
	speciesInfo := SpeciesInfo{}

	req, err := http.NewRequest("GET", baseURL+"species", nil)
	if err != nil {
		return speciesInfo, errors.Wrap(err, "unable to create request")
	}

	q := req.URL.Query()
	q.Add("Genus", genus)
	q.Add("Species", species)

	req.URL.RawQuery = q.Encode()

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return speciesInfo, errors.Wrap(err, "unable to query fishbase")
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return speciesInfo, errors.Wrap(err, "unable to read response body")
	}

	var si SpeciesInfo
	if err := json.Unmarshal(body, &si); err != nil {
		return speciesInfo, errors.Wrap(err, "unable to unmarshal response")
	}

	if len(si.Data) == 0 {
		return speciesInfo, fmt.Errorf("Genus '%s' and Species '%s' returned no results", genus, species)
	}

	return si, nil
}
