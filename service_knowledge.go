package splunk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
)

// KnowledgeService is an interface for interacting with the splunk knowledge API
type KnowledgeService interface {
	ListExtractions(opts *ListOptions) ([]map[string]interface{}, error)
	ListTransforms(opts *ListOptions) ([]map[string]interface{}, error)
	CreateExtraction(name, stanza, extractionType, value string, acl *ACL) error
	DeleteExtraction(stanza, extractionType, value string) error
}

type knowledgeService service

// TODO: pagination
func (ks *knowledgeService) ListExtractions(opts *ListOptions) ([]map[string]interface{}, error) {
	data, err := ks.client.NewRequest("GET", "/services/data/props/extractions", nil)
	if err != nil {
		return nil, err
	}

	var extractions []map[string]interface{}
	err = json.Unmarshal(data.Entry, &extractions)
	if err != nil {
		return nil, err
	}

	return extractions, nil
}

// TODO: pagination
func (ks *knowledgeService) ListTransforms(opts *ListOptions) ([]map[string]interface{}, error) {
	data, err := ks.client.NewRequest("GET", "/services/data/transforms/extractions", nil)
	if err != nil {
		return nil, err
	}

	var transforms []map[string]interface{}
	err = json.Unmarshal(data.Entry, &transforms)
	if err != nil {
		return nil, err
	}

	return transforms, nil
}

func (ks *knowledgeService) CreateExtraction(name, stanza, extractionType, value string, acl *ACL) error {
	form := url.Values{}
	form.Set("name", name)
	form.Set("stanza", stanza)
	form.Set("type", extractionType)
	form.Set("value", value)

	_, err := ks.client.NewRequest("POST", "/services/data/props/extractions", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	if acl != nil {
		uri := path.Join(encodedExtractionURI(stanza, extractionType, value), "acl")

		_, err := ks.client.NewRequest("POST", uri, strings.NewReader(acl.Encode()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ks *knowledgeService) DeleteExtraction(stanza, extractionType, value string) error {
	_, err := ks.client.NewRequest("DELETE", encodedExtractionURI(stanza, extractionType, value), nil)
	if err != nil {
		return err
	}

	return nil
}

func encodedExtractionURI(stanza, extractionType, value string) string {
	return path.Join(
		"/services/data/props/extractions",
		url.PathEscape(fmt.Sprintf("%s : %s-%s", stanza, extractionType, value)),
	)
}
