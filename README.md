[![PkgGoDev](https://pkg.go.dev/badge/github.com/tommyblue/smugmug-backup)](https://pkg.go.dev/github.com/tommyblue/smugmug-backup)

# SmugMug backup

![Go](https://github.com/tommyblue/smugmug-backup/workflows/Go/badge.svg)

Makes a full backup of a [SmugMug](https://www.smugmug.com/) account (images and videos
are supported).

The program loops over the images and videos of the user's albums and saves them in the destination
folder, replicating the SmugMug paths.

You can run the app multiple times, all exising files will be skipped if their sizes match.

## Releases

Releases for multiple systems can be found in the
[project releases page](https://github.com/tommyblue/smugmug-backup/releases)

## Configuration

The app expects to find the [TOML](https://github.com/toml-lang/toml) configuration file in
`./config.toml` or `$HOME/.smgmg/config.toml`.

The configuration file must have this structure:

```toml
[authentication]
username = "<SmugMug username>"
api_key = "<API Key>"
api_secret = "<API Secret>"
user_token = "<User Token>"
user_secret = "<User Secret>"

[store]
destination = "<Backup destination folder>"
file_names = "<Filename with template replacements>"
```

All values can be overridden by environment variables, that have the following names:

```sh
SMGMG_BK_USERNAME = "<SmugMug username>"
SMGMG_BK_API_KEY = "<API Key>"
SMGMG_BK_API_SECRET = "<API Secret>"
SMGMG_BK_USER_TOKEN = "<User Token>"
SMGMG_BK_USER_SECRET = "<User Secret>"
SMGMG_BK_DESTINATION = "<Backup destination folder>"
SMGMG_BK_FILE_NAMES = "<Filename with template replacements>"
```

All configuration values are required. They can be omitted in the configuration file
as long as they are overridden by environment values.

The **username** can be found in the first part of the url in your SmugMug's homepage.  
In my case the url is [https://tommyblue.smugmug.com/](https://tommyblue.smugmug.com/) and the
username is `tommyblue`.

The **destination** is the local path to save SmugMug pictures and videos.  
If the folder is not empty, then only new or changed files will be downloaded.

**file_names** is a string including template replacements that will be used to build the file
names for the files on disk. Accepted keys are `FileName`, `ImageKey`, `ArchivedMD5` and `UploadKey`
and their values comes from the AlbumImage API response. If an invalid replacement is used,
an error is returned. If the conf key is omitted or is empty, then `{{.FileName}}` is used.  

**api_key**, **api_secret**, **user_token** and **user_secret** are the required credentials for
authenticating with the SmugMug API.  
See [credentials](#credentials) below for details about how to obtain them.

## Run

Once the configuration file and/or the environment variables are set,
you can perform the account backup with:

```sh
./smugmug-backup
```

Running the backup can take a lot of time, depending on the size of your account and the
connection speed. Check the command line logs to see what's going on.

## Credentials

SmugMug requires *OAuth1 authentication*. OAuth1 requires 4 values: an API key and secret that
identify the app with SmugMug and that you can get from SmugMug registering the app and
user credentials (token and secret) that each user must obtain (those credentials authorize
the app to access the user's data).

### Obtain API keys

Apply for an API key from
[https://api.smugmug.com/api/developer/apply](https://api.smugmug.com/api/developer/apply)
and wait the app to be authorized.

Once you have obtained the API key and secret, save them in the `authentication` section
of the configuration file:

```toml
[authentication]
api_key = "<API Key>"
api_secret = "<API Secret>"
```

### Obtain Tokens

Once your app has been accepted by SmugMug and you got the API key and secret, then go to your
[Account Settings > Privacy page](https://www.smugmug.com/app/account/settings/?#section=privacy)
and scroll down to "Authorized Services", where you'll find the app and a link to see the tokens.

Add them to the `authentication` section of the configuration file:

```toml
[authentication]
user_token = "<User Token>"
user_secret = "<User Secret>"
```

#### Alternative ways to obtain the tokens

Based on the examples from SmugMug (that you can find in the [get_tokens](./get_tokens) folder,
I've written a small web app that can help everyone to obtain their user token and secret.

The app has its own [GitHub repo](https://github.com/tommyblue/smugmug-api-authenticator) and a live
version is deployed to heroku at
[https://smugmug-api-authenticator.herokuapp.com/](https://smugmug-api-authenticator.herokuapp.com/).  
You can use that app, it doesn't store any personal data in the server, but (as you should) you
don't trust me, you can easily clone the
[GitHub repo](https://github.com/tommyblue/smugmug-api-authenticator), check the code, run the
app locally and get the tokens.

If you prefer to use the console, the `get_tokens` folder contains a script from SmugMug to obtain
the OAuth1 tokens.
You need to create a `config.json` file with your API key/secret using `example.json` as example.
Then, using a python3 environment, run `run-console.sh`.
The script will show you a link you must open with your browser. SmugMug will give you a 6-digit
code you must then paste to the console prompt.
That's the last step, the console will show the user token and secret

## Build and install

To build and install the program from source:

```sh
git clone https://github.com/tommyblue/smugmug-backup.git
cd smugmug-backup
make build
```

This will produce the `./smugmug-backup` binary.  
More `make` commands are available, run `make help` to get help

## Debug for errors

To increase the logging, export a `DEBUG=1` environment variable:

```sh
DEBUG=1 ./smugmug-backup
```

## Credits

OAuth1 signature has been heavily inspired by
[https://github.com/gomodule/oauth1](https://github.com/gomodule/oauth1)

The code in the `get_tokens` folder is a copy of
[https://gist.github.com/smugkarl/10046914](https://gist.github.com/smugkarl/10046914)

## Bugs and contributing

[Package documentation](https://pkg.go.dev/github.com/tommyblue/smugmug-backup?tab=doc)

If you find a bug or want to suggest something, please
[open an issue](https://github.com/tommyblue/smugmug-backup/issues/new).

If you want to contribute to this project, fork the repo and open a pull-request.  
Contributing is more than welcome :smile:
