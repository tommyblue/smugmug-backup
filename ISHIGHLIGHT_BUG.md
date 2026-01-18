# IsHighlight Bug Documentation

## Problem Description

The `IsHighlight` field in `images_metadata.csv` is not correctly identifying the actual cover/highlight image of SmugMug albums.

## Current Behavior

The current implementation attempts to call the SmugMug API endpoint `/api/v2/highlight/node/{nodeId}` to retrieve the ImageKey of the album's cover image. However, the images marked as `IsHighlight=true` do NOT match the actual cover images shown in SmugMug.

### Example Cases

**Album: Corno alle Scale Luglio 2025**
- Expected cover ImageKey: `Hvxb9RD` (as shown in SmugMug: https://www.smugmug.com/app/organize/Album/Corno-alle-Scale-Luglio-2025/i-Hvxb9RD)
- Actual marked in CSV: `tvB2vwd` (incorrect)

**Album: Bagni di Mario 2024**
- Expected cover ImageKey: `wPZmTgm` (as shown in SmugMug: https://www.lorenzoperone.it/Album/2024-Bagni-di-Mario/n-nrHn95/i-wPZmTgm/A)
- Actual marked in CSV: `3dZmzbc` (incorrect)

## Root Cause

The SmugMug album API returns `Uris.HighlightImage.Uri` pointing to `/api/v2/highlight/node/{nodeId}` instead of directly to `/api/v2/image/{ImageKey}`.

When calling the highlight node endpoint, the response structure is **unknown/undocumented**, and our current parsing logic in `highlightImageKey()` function does not correctly extract the actual cover image's ImageKey.

## Current Implementation

### Code Location: `api.go`

```go
func (w *Worker) highlightImageKey(nodeId string) string {
	if nodeId == "" {
		return ""
	}
	uri := fmt.Sprintf("/api/v2/highlight/node/%s", nodeId)
	var resp highlightNodeResponse
	if err := w.req.get(uri, &resp); err != nil {
		log.Debugf("error getting highlight image for node %s: %v", nodeId, err)
		return ""
	}

	// Try to extract ImageKey from various possible response structures
	// Format 1: Direct Image.ImageKey
	if resp.Response.Image.ImageKey != "" {
		log.Debugf("Found highlight ImageKey via Image.ImageKey: %s", resp.Response.Image.ImageKey)
		return resp.Response.Image.ImageKey
	}
	// ... (tries 6 different formats)
}
```

### Current Fallback Behavior

If the API call fails or returns an empty ImageKey, **no images are marked as highlight** (`IsHighlight=false` for all images in the album). This is a conservative approach to avoid marking incorrect images.

## Investigation Needed

To fix this issue, we need to:

1. **Log the raw JSON response** from `/api/v2/highlight/node/{nodeId}` to understand the actual response structure
2. **Identify which field** in the response contains the correct ImageKey
3. **Update the `highlightNodeResponse` struct** and parsing logic accordingly

### Proposed Debug Approach

Add debug logging to print the complete API response:

```go
func (w *Worker) highlightImageKey(nodeId string) string {
	if nodeId == "" {
		return ""
	}
	uri := fmt.Sprintf("/api/v2/highlight/node/%s", nodeId)
	
	// TODO: Add raw response logging here
	var rawResponse map[string]interface{}
	if err := w.req.get(uri, &rawResponse); err != nil {
		log.Debugf("error getting highlight image for node %s: %v", nodeId, err)
		return ""
	}
	log.Debugf("Raw highlight node response for %s: %+v", nodeId, rawResponse)
	
	// Continue with structured parsing...
}
```

## Impact

- Users cannot reliably identify which image is the actual album cover
- The `IsHighlight` field in CSV exports is unreliable
- May result in incorrect cover images when migrating from SmugMug

## Status

**OPEN** - Investigation required to understand SmugMug API response structure for highlight nodes.

## Related Files

- `api.go` - Contains `highlightImageKey()` function
- `json_structs.go` - Contains `highlightNodeResponse` struct
- `csv.go` - Uses highlight ImageKey to mark images in CSV

## Testing

All 19 unit tests pass, but they use mocked API responses. Real-world API behavior differs from the mock.

**Test file:** `csv_test.go`
```go
type mockHighlightHandler struct {
	highlightImageKey string
}

func (m *mockHighlightHandler) get(url string, obj interface{}) error {
	if resp, ok := obj.(*highlightNodeResponse); ok {
		resp.Response.Image.ImageKey = m.highlightImageKey
		return nil
	}
	return nil
}
```

The mock always returns a simple `Image.ImageKey` field, but the real API may have a different structure.
