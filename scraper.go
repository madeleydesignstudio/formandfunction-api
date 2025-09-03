package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const graphQLURL = "https://www.travisperkins.co.uk/graphql?op=tpplcProductCollectionAvailability"

// GraphQLRequest represents the payload for the GraphQL request.
type GraphQLRequest struct {
	OperationName string    `json:"operationName"`
	Query         string    `json:"query"`
	Variables     Variables `json:"variables"`
}

// Variables represents the variables for the GraphQL query.
type Variables struct {
	ProductID string `json:"productId"`
	Postcode  string `json:"postcode"`
	BrandID   string `json:"brandId"`
}

// GraphQLResponse represents the expected response from the GraphQL API.
type GraphQLResponse struct {
	Data struct {
		TpplcBrand struct {
			ProductCollectionAvailability []struct {
				BranchID   string  `json:"branchId"`
				StockLevel float64 `json:"stockLevel"`
			} `json:"productCollectionAvailability"`
		} `json:"tpplcBrand"`
	} `json:"data"`
}

// GetStockStatus sends a request to the GraphQL API to get stock info.
func GetStockStatus(productID, postcode string) (string, error) {
	requestBody := GraphQLRequest{
		OperationName: "tpplcProductCollectionAvailability",
		Query: `query tpplcProductCollectionAvailability($branchId: String, $branchLimit: Int, $postcode: String, $productId: String!, $withinRadius: Float, $brandId: ID!) {\n  tpplcBrand(brandId: $brandId) {\n    productCollectionAvailability(\n      branchId: $branchId
      branchLimit: $branchLimit
      postcode: $postcode
      productId: $productId
      withinRadius: $withinRadius
    ) {\n      branchId
      stockLevel
      stockUom
      __typename
    }
    __typename
  }
}`,
		Variables: Variables{
			ProductID: productID,
			Postcode:  postcode,
			BrandID:   "tp",
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal graphql request: %w", err)
	}

	req, err := http.NewRequest("POST", graphQLURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to graphql api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("graphql api request failed with status: %s", resp.Status)
	}

	var gqlResponse GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResponse); err != nil {
		return "", fmt.Errorf("failed to decode graphql response: %w", err)
	}

	if len(gqlResponse.Data.TpplcBrand.ProductCollectionAvailability) > 0 {
		for _, branch := range gqlResponse.Data.TpplcBrand.ProductCollectionAvailability {
			if branch.StockLevel > 0 {
				return "InStock", nil
			}
		}
		return "OutOfStock", nil
	}

	return "NotAvailable", nil
}
