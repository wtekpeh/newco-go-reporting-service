package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type InternetSearchService struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

type TavilySearchRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type TavilySearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type TavilySearchResponse struct {
	Query   string               `json:"query"`
	Answer  string               `json:"answer"`
	Results []TavilySearchResult `json:"results"`
}

func NewInternetSearchService(
	baseURL string,
	apiKey string,
) *InternetSearchService {

	return &InternetSearchService{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{},
	}
}

func (s *InternetSearchService) Search(
	ctx context.Context,
	query string,
) (TavilySearchResponse, error) {

	requestBody := TavilySearchRequest{
		APIKey:      s.APIKey,
		Query:       query,
		SearchDepth: "basic",
		MaxResults:  5,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return TavilySearchResponse{}, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.BaseURL,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return TavilySearchResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := s.Client.Do(request)
	if err != nil {
		return TavilySearchResponse{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return TavilySearchResponse{}, err
	}

	if response.StatusCode != http.StatusOK {
		return TavilySearchResponse{}, fmt.Errorf(
			"tavily returned status %d: %s",
			response.StatusCode,
			string(responseBody),
		)
	}

	var searchResponse TavilySearchResponse

	err = json.Unmarshal(responseBody, &searchResponse)
	if err != nil {
		return TavilySearchResponse{}, err
	}

	return searchResponse, nil
}
