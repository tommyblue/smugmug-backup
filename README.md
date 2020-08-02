# SmugMug backup

![Go](https://github.com/tommyblue/smugmug-backup/workflows/Go/badge.svg)

Makes a full backup of a [SmugMug](https://www.smugmug.com/) account (images and videos are supported).

The program loops over the images and videos of the user's albums and saves them in the destination
folder, replicating the SmugMug paths.

You can run the app multiple times, all exising files will be skipped if their sizes match.

## Releases

Releases for multiple systems can be found in the [project releases page](https://github.com/tommyblue/smugmug-backup/releases)

## Credentials

SmugMug requires OAuth1 authentication. OAuth1 requires 4 values: API key and secret you can get
from SmugMug registering the app (they're generic for the app) and user token and secret each user
must obtain (they authorize the app to access the user's data).

### Obtain API keys

Apply for an API key: [https://api.smugmug.com/api/developer/apply](https://api.smugmug.com/api/developer/apply)

You'll get an API key and secret, save them as environment variables:

```sh
export API_KEY="<key>"
export API_SECRET="<secret>"
```

### Obtain Tokens

Once your app has been accepted by SmugMug and you got the API key and secret, then go to your [Account Settings > Privacy page](https://www.smugmug.com/app/account/settings/?#section=privacy) and scroll down to "Authorized Services", where you'll find the app and a link to see the tokens.  
They must be exported as environment variables:

```sh
export USER_TOKEN="<Access Token>"
export USER_SECRET="<Token Secret>"
```

#### Alternative ways to obtain the tokens

Based on the examples from SmugMug (that you can find in the [get_tokens](./get_tokens) folder,
I've written a small web app that can help everyone to obtain their user token and secret.

The app has its own [GitHub repo](https://github.com/tommyblue/smugmug-api-authenticator) and a live
version is deployed to heroku at [https://smugmug-api-authenticator.herokuapp.com/](https://smugmug-api-authenticator.herokuapp.com/).  
You can use that app, it doesn't store any personal data in the server, but (as you should) you don't
trust me, you can easily clone the [GitHub repo](https://github.com/tommyblue/smugmug-api-authenticator),
check the code, run the app locally and get the tokens.

If you prefer to use the console, the `get_tokens` folder contains a script from SmugMug to obtain
the OAuth1 tokens.
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
git clone https://github.com/tommyblue/smugmug-backup.git
make build
```

More `make` commands are available, run `make help` to get help

## Run

With all the four environment variables set, you can run the program with:

```sh
./smugmug-backup -user <username> -destination <path of downloads>
```

### Command line options

#### -user \<username\>

The username can be found in the first part of the url in your SmugMug's homepage.  
In my case the url is [https://tommyblue.smugmug.com/](https://tommyblue.smugmug.com/) and the username is `tommyblue`.

#### -destination \<path\>

Local path to save SmugMug pictures and videos. If not empty, only new or changed files will be downloaded.

## Debug for errors

To increase the logging, export a `DEBUG=1` environment variable:

```sh
DEBUG=1 ./smugmug-backup -user <username> -destination <path of downloads>
```

## Credits

OAuth1 signature has been heavily inspired by [https://github.com/gomodule/oauth1](https://github.com/gomodule/oauth1)

The code in the `get_tokens` folder is a copy of [https://gist.github.com/smugkarl/10046914](https://gist.github.com/smugkarl/10046914)

## Bugs and contributing

If you find a bug or want to suggest something, please [open an issue](https://github.com/tommyblue/smugmug-backup/issues/new).

If you want to contribute to this project, fork the repo and open a pull-request. Contributing is more than welcome :smile:
