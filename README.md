# smugmug-backup
Makes a full backup of a SmugMug account.

OAuth1 signature has been heavily inspired by https://github.com/gomodule/oauth1

## Setup and run

```
export API_KEY=""
export API_SECRET=""
export USER_TOKEN=""
export USER_SECRET=""
go run *.go -user <username> -destination <path of downloads>
```
