## Description

This PR implements **comprehensive SmugMug data export** with complete metadata for albums and images/videos, transforming the minimal CSV export into a full-featured data archival system.

## 🎯 What's New

### Complete CSV Metadata Export

**Original**: Single `metadata.csv` with 7 basic fields:
- Filename, Type, ArchivedUri, Caption, Keywords, Latitude, Longitude

**New**: Two separate CSV files with **51 total fields** of complete SmugMug data:

#### `albums_metadata.csv` (19 fields)
Complete album information:
- AlbumKey, URLPath, Title, Name, NiceName, Description, Keywords
- Date, LastUpdated, ImagesLastUpdated, ImageCount
- Privacy, SecurityType, SortMethod, SortDirection
- WebUri, AllowDownloads, PasswordHint, Protected

#### `images_metadata.csv` (32 fields)
Complete image/video metadata:
- **Identity**: Filename, Type, AlbumKey, AlbumPath, ImageKey, Title, Caption, Keywords
- **Technical**: Format, Width, Height, OriginalWidth, OriginalHeight, Size, ArchivedSize, ArchivedUri, ArchivedMD5, UploadKey
- **GPS**: Latitude, Longitude, Altitude
- **Dates**: DateTimeOriginal, DateTimeUploaded
- **Flags**: Hidden, Watermarked, Collectable, IsArchive, IsHighlight (cover image)
- **Status**: Status, SubStatus, WebUri, ThumbnailUrl

### Key Features

1. **Dual-CSV architecture** - Separate albums and images to avoid data duplication, with AlbumKey for correlation
2. **Album metadata propagation** - Album fields (Title, Description, Keywords, dates) included in each image record for standalone usability
3. **Cover image identification** - `IsHighlight` field correctly marks the album cover image by calling SmugMug API `/api/v2/highlight/node/{nodeId}`
4. **Complete data portability** - Perfect for users migrating from SmugMug or creating comprehensive local archives
5. **Expanded data structures**:
   - `album` struct: 3 → 20 fields (captures all API data)
   - `albumImage` struct: 18 → 38 fields (complete metadata)

### Additional Improvements

- Fixed Altitude field type from string to int (matches SmugMug API)
- Fixed HighlightImageUri extraction from nested Uris structure
- Comprehensive test coverage (19 tests passing)

## Type of change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [x] New feature (non-breaking change which adds functionality)
- [x] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [x] This change requires documentation update

## ⚠️ Breaking Changes

**CSV export users**: The old `metadata.csv` (7 fields) is replaced by two new files:
- `albums_metadata.csv` (19 fields)
- `images_metadata.csv` (32 fields)

Users with existing automation or scripts parsing `metadata.csv` will need to update their code to use the new files and field names.

## How Has This Been Tested?

1. **Full backup**: 386 albums, 38,195 images successfully processed and verified
2. **CSV validation**: All 51 fields across both files populated correctly
3. **Cover image verification**: Compared `IsHighlight` images with SmugMug web interface - 100% match on sampled albums
4. **Unit tests**: All 19 tests passing (`go test ./...`)

### Verified Examples (Cover Images):

| Album | Expected ImageKey | Verified |
|-------|------------------|----------|
| Bagni di Mario 2024 | `wPZmTgm` | ✅ |
| Corno alle Scale | `Hvxb9RD` | ✅ |
| Navibulgar | `S7tXRZc` | ✅ |

## Checklist

- [x] My code follows the style guidelines of this project
- [x] I have performed a self-review of my own code
- [x] I have commented my code, particularly in hard-to-understand areas
- [x] I have made corresponding changes to the documentation (README.md, CHANGELOG.md)
- [x] My changes generate no new warnings
- [x] I have checked my code and corrected any misspellings
- [x] All new and existing tests passed

## Additional Notes

- The file `ISHIGHLIGHT_BUG.md` documents the development process for implementing cover image identification and can be kept for historical reference or removed
- This implementation maintains backward compatibility for non-CSV users (image/video download logic unchanged)
- Configuration option `store.write_csv` continues to work (now exports to two files instead of one)
