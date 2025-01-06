# journal

```bash
go mod tidy

JOURNAL_OLLAMA_URL=http://localhost:11434 \
JOURNAL_OLLAMA_INSECURE=1 \
JOURNAL_SESSION_KEY=zsdftyhghuijk345e6r7t8y9hio \
JOURNAL_PASSWORD=abcd1234 \
JOURNAL_SESSION_SECURE=0 \
go run main.go
```

## Setting up Ollama

cwpearson/journal uses "llama3.2:3b" which wants 3.4 GiB of RAM available

* `JOURNAL_OLLAMA_URL`
* `JOURNAL_OLLAMA_INSECURE`

## Docker-Compose Example

```bash
docker run --rm -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

```yaml

```

## Configuration Options

You must provide the following

* `JOURNAL_SESSION_KEY`: the key used to secure the cookie session.
* `JOURNAL_PASSWORD`: the password needed to authenticate in the browser.

The following configuration options are optional

* `JOURNAL_PORT`: what port for the server to listen on (default: `8080`)
* `JOURNAL_SESSION_SECURE`: set to `0` if you are accessing the app without HTTPS (default `1`)
* `JOURNAL_OLLAMA_URL`: where to access Ollama enpoints
* `JOURNAL_OLLAMA_INSECURE`: set to `1` if access Ollama over http

## Roadmap

- Footer
  - [ ] Build SHA
  - [ ] Build date
- Images
- Config
  - [ ] `JOURNAL_SITE_NAME`
  - [x] `JOURNAL_PASSWORD`
  - [x] `JOURNAL_OLLAMA_URL`
  - [x] `JOURNAL_OLLAMA_INSECURE`
  - [x] `JOURNAL_PORT`
- Tags
  - [ ] Clicking on tags on history page
  - [ ] Retrieve entries by tag
- Docker
  - [x] ghcr.io publish
- Ollama
  - [ ] better handling of missing / errored Ollama

## Setting up GHCR Publish

* github > settings > developer settings > personal access tokens > access tokens (classic) > generate new token (classic)
  * `write:packages`
* copy the token
* set as `GHCR_TOKEN` in actions secrets on this repo