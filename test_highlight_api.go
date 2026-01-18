package smugmug

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// TestHighlightNodeStructure attempts to call a real SmugMug API endpoint
// and dump the response to understand the structure of /api/v2/highlight/node/{nodeId}
func TestHighlightNodeStructure(t *testing.T) {
	// Example node IDs from real albums
	nodeIds := []string{"fpT33v", "nrHn95", "vXKKJB"}

	// Note: This test requires valid OAuth credentials set in config
	// It's disabled by default - uncomment to debug
	t.Skip("Skipping real API test - uncomment to debug highlight node structure")

	for _, nodeId := range nodeIds {
		url := fmt.Sprintf("https://api.smugmug.com/api/v2/highlight/node/%s", nodeId)

		resp, err := http.Get(url)
		if err != nil {
			t.Logf("Error calling %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Logf("Error reading body: %v", err)
			continue
		}

		t.Logf("\n=== Node %s ===", nodeId)
		t.Logf("Status: %d", resp.StatusCode)
		t.Logf("Raw Response:\n%s", string(body))

		// Try to unmarshal with generic structure
		var generic map[string]interface{}
		if err := json.Unmarshal(body, &generic); err != nil {
			t.Logf("Error parsing JSON: %v", err)
			continue
		}

		// Pretty print
		prettyJson, _ := json.MarshalIndent(generic, "", "  ")
		t.Logf("Pretty JSON:\n%s", string(prettyJson))
	}
}

// TestHighlightNodeResponseStructures tests various possible response structures
// based on SmugMug API patterns
func TestHighlightNodeResponseStructures(t *testing.T) {
	// Based on SmugMug API patterns, highlight nodes might contain:

	// Possibility 1: Direct image reference
	type possibility1 struct {
		Response struct {
			Image struct {
				ImageKey string `json:"ImageKey"`
				URI      string `json:"Uri"`
			} `json:"Image"`
			Uris struct {
				Image struct {
					URI string `json:"Uri"`
				} `json:"Image"`
			} `json:"Uris"`
		} `json:"Response"`
	}

	// Possibility 2: Highlight image nested
	type possibility2 struct {
		Response struct {
			HighlightImage struct {
				ImageKey string `json:"ImageKey"`
				URI      string `json:"Uri"`
			} `json:"HighlightImage"`
			Uris struct {
				HighlightImage struct {
					URI string `json:"Uri"`
				} `json:"HighlightImage"`
			} `json:"Uris"`
		} `json:"Response"`
	}

	// Possibility 3: Node contains Image Uris
	type possibility3 struct {
		Response struct {
			Uris struct {
				Image struct {
					URI string `json:"Uri"`
				} `json:"Image"`
				AlbumImage struct {
					URI string `json:"Uri"`
				} `json:"AlbumImage"`
			} `json:"Uris"`
		} `json:"Response"`
	}

	t.Logf("Possible structures for highlight node response:")
	t.Logf("1. Direct Image.ImageKey: %T", possibility1{})
	t.Logf("2. HighlightImage.ImageKey: %T", possibility2{})
	t.Logf("3. Uris.Image URI (needs parsing): %T", possibility3{})
	t.Logf("\nThe actual structure should be determined by calling the real API")
}
