# SmugMug backup

![Go](https://github.com/tommyblue/smugmug-backup/workflows/Go/badge.svg)

Makes a full backup of a [SmugMug](https://www.smugmug.com/) account (images and videos are supported).

The program loops over the images and videos of the user's albums and saves them in the destination
folder, replicating the SmugMug paths.

You can run the app multiple times, all exising files will be skipped if their sizes match.

**Note on multiple runs**

The app skips existing files but doesn't check their sizes. So if you interrupt a running execution
of the program while it's saving a file, it will exist but will also probably be corrupted.

I'd like to manage this situation, but in the meanwhile when you interrupt an execution, take notes
of the last files shown in the console and check if they're valid, deleting them otherwise.

## Credentials

SmugMug requires OAuth1 authentication. OAuth1 requires 4 values: API key and secret you can get
from SmugMug registering the app (they're generic for the app) and user token and secret each user
must obtain (they authorize the app to access the user's data).

### Obtain API keys

Apply for an API key: https://api.smugmug.com/api/developer/apply

You'll get an API key and secret, save them as environment variables:

```sh
export API_KEY="<key>"
export API_SECRET="<secret>"
```

### Obtain Tokens

The `get_tokens` folder contains a script from SmugMug to obtain OAuth1 tokens.
You need to create a `config.json` file with your API key/secret using `example.json` as example.
Then, using a python3 environment, run `run-console.sh`.
The script will show you a link you must open with your browser. SmugMug will give you a 6-digit
code you must then paste to the console prompt.
That's the last step, the console will show the user token and secret. Export them:

```sh
export USER_TOKEN="<token>"
export USER_SECRET="<secret>"
```

## Build and install

To build and install the program:

```sh
go get github.com/tommyblue/smugmug-backup
cd $GOPATH/src/github.com/tommyblue/smugmug-backup
go get ./...
go install
```

## Run

With all the four environment variables set, you can run the program with:

```sh
$GOPATH/bin/smugmug-backup -user <username> -destination <path of downloads>
```

The **username** can be found in the first part of the url in your SmugMug's homepage.
In my case the url is https://tommyblue.smugmug.com/ and the username `tommyblue`

I suggest adding `$GOPATH/bin` to the `$PATH` so you can avoid writing the full path of the program.

## Run without install

```
go run *.go -user <username> -destination <path of downloads>
```

## Debug for errors

To increase the logging, export a `DEBUG=1` environment variable:

```sh
DEBUG=1 $GOPATH/bin/smugmug-backup -user <username> -destination <path of downloads>
```

## Credits

OAuth1 signature has been heavily inspired by https://github.com/gomodule/oauth1

The code in the `get_tokens` folder is a copy of https://gist.github.com/smugkarl/10046914
