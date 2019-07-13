package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"net/http"
)

type response struct {
	Items [] item `yaml:"items"`
}

type item struct {
	Spec Spec `yaml:"spec"`
}

type Spec struct {
	Containers []corev1.Container `yaml:"containers"`
	Selector   Selector           `yaml:"selector"`
}

type Selector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type WebClient struct {
	Url    string
	Token  string
	Client *http.Client
}

func (webClient *WebClient) LoadSidecarConfig() ([]Spec, error) {
	url := webClient.Url
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+webClient.Token)

	r, err := webClient.Client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var resp response
	if err := yaml.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	var list [] Spec
	for _, item := range resp.Items {
		list = append(list, item.Spec)
	}
	return list, nil
}

