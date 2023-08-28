# SmugMug backup

[![PkgGoDev](https://pkg.go.dev/badge/github.com/tommyblue/smugmug-backup)](https://pkg.go.dev/github.com/tommyblue/smugmug-backup) ![Go](https://github.com/tommyblue/smugmug-backup/workflows/Go/badge.svg)

Makes a full backup of a [SmugMug](https://www.smugmug.com/) account (images and videos
are supported).

The program loops over the images and videos of the user's albums and saves them in the destination
folder, replicating the SmugMug paths.

You can run the app multiple times, all exising files will be skipped if their sizes match.

- [SmugMug backup](#smugmug-backup)
  - [Releases](#releases)
  - [Configuration](#configuration)
  - [Run](#run)
  - [Credentials](#credentials)
    - [Obtain API keys](#obtain-api-keys)
    - [Obtain Tokens](#obtain-tokens)
      - [Alternative ways to obtain the tokens](#alternative-ways-to-obtain-the-tokens)
  - [Build and install](#build-and-install)
  - [Debug for errors](#debug-for-errors)
  - [Visualize software stats](#visualize-software-stats)
  - [Credits](#credits)
  - [Bugs and contributing](#bugs-and-contributing)

## Releases

Releases for multiple systems (including ARM) can be found in the
[project releases page](https://github.com/tommyblue/smugmug-backup/releases)

## Configuration

The app expects to find the [TOML](https://github.com/toml-lang/toml) configuration file in
`./config.toml` or `$HOME/.smgmg/config.toml`.

For the configuration structure, check the [config.example.toml](./config.example.toml) file.

Some values can be overridden by environment variables, that have the following names:

```sh
SMGMG_BK_API_KEY = "<API Key>"
SMGMG_BK_API_SECRET = "<API Secret>"
SMGMG_BK_USER_TOKEN = "<User Token>"
SMGMG_BK_USER_SECRET = "<User Secret>"
SMGMG_BK_DESTINATION = "<Backup destination folder>"
SMGMG_BK_FILE_NAMES = "<Filename with template replacements>"
```

| Configuration identifier   | Required | Default Value   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
|----------------------------|----------|-----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| authentication.*           | Yes      |                 | See [credentials](#credentials) below for details about how to obtain them.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| store.destination          | Yes      |                 | Local path to save SmugMug pictures and videos into. If the folder is not empty, then only new or changed files will be downloaded.    **Windows users**  The value of `destination` must use slash `/` or double backslash `\\`  Examples:   ```toml  destination = "C:/folder/subfolder"  destination = "C:\\folder\\subfolder"  destination = "/folder/subfolder" # This writes to the primary partition C: ```                                                                                                                                                                                                                                                                                                                                 |
| store.file_names           | No       | `{{.FileName}}` | Is a string including template replacements that will be used to build the file names for the files on disk. Accepted keys are `FileName`, `ImageKey`, `ArchivedMD5` and `UploadKey` and their values comes from the AlbumImage API response. If an invalid replacement is used, an error is returned                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| store.use_metadata_times   | No       | false           | When true, the last modification timestamp of the objects will be set based on SmugMug metadata for newly downloaded files. If also **force_metadata_times** is true, then the timestamp is applied to all existing files. This configuration can be required if you notice that the images creation datetime is wrong by ~7h. This is a bug in the SmugMug Uploader: "Our uploader process currently isn't time zone aware and takes the DateTimeOriginal field without time zone information".  The solution is to use the Metadata API endpoint to retrieve the EXIF informations, but it requires an additional API call for each image/video.    In my case, a full backup that requires ~10 minutes, increases to 2+ hours with this option. |
| store.force_metadata_times | No       | false           |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| store.write_csv            | No       | false           | When true, a `metadata.csv` file is created (or overwritten) storing some information about the user files (both downloaded or skipped).                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| store.force_video_download | No       | false           | When true, videos are downloaded also if marked as "under processing"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| store.concurrent_albums    | No       | 1               | Number of concurrently analyzed albums. It multiplies API calls, don't stress this number or you'll be rate limited. 5 is a good start.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| store.concurrent_downloads | No       | 1               | Number of concurrently downloaded images and videos.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |

## Run

Once the configuration file and/or the environment variables are set,
you can perform the account backup with:

```sh
./smugmug-backup
```

Running the backup can take a lot of time, depending on the size of your account and the
connection speed. Check the command line logs to see what's going on.

## Credentials

SmugMug requires _OAuth1 authentication_. OAuth1 requires 4 values: an API key and secret that
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

## Visualize software stats

In case you want to monitor the software performances you can use the `-stats` flag to enable
[statsviz](https://github.com/arl/statsviz) debug page (available at http://localhost:6060/debug/statsviz/)

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
