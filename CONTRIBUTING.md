# Contributing guidelines

This project started as a personal tool to backup my account, but over time I realized
a lot of people had the same necessity. So, first off, thank you for considering contributing
to Smugmug backup.

Before interacting with the community around this project, be sure to read the
[code of conduct](./code_of_conduct.md).

## Did you find a potential bug? Suggestions?

If you find a bug or want to suggest something, please
[open an issue](https://github.com/tommyblue/smugmug-backup/issues/new). I'll try to reply
as soon as possibile. Note that I generally use Linux or Mac, so Windows-related bugs
happen more often as I'm not actively using that o.s. Anyway, I'll do my best to reproduce
the bug and fix it!

## Want to contribute yourself?

If you're considering to add your contribution, you're more than welcome!

Feel free to open a new pull request to add code, docs, translations, whatever.

Contributing is more than welcome :smile:

I'm dropping here the link to
[package documentation](https://pkg.go.dev/github.com/tommyblue/smugmug-backup?tab=doc), in
case you need that.

### Mocking the API server

This package contains a mock server that can be executed with `go run ./cmd/api_mock` so that you
don't need to have a real SmugMug account. The server will listen on http://127.0.0.1:3000  
To use it run the smugmug-backup binary with the `-mock` flag.
