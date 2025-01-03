# journal

```bash
go mod tidy

JOURNAL_OLLAMA_URL=http://localhost:11434 \
JOURNAL_OLLAMA_INSECURE=1 \
go run main.go
```

## Setting up Ollama

cwpearson/journal uses "llama3.2:3b"

* `JOURNAL_OLLAMA_URL`
* `JOURNAL_OLLAMA_INSECURE`

## Docker-Compose Example

```bash
docker run --rm -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

```yaml

```

## Roadmap

- Images
- Config
  - [ ] `JOURNAL_SITE_NAME`
  - [ ] `JOURNAL_PASSWORD`
  - [x] `JOURNAL_OLLAMA_URL`
  - [x] `JOURNAL_OLLAMA_INSECURE`
- Docker
  - [x] ghcr.io publish

## Setting up GHCR Publish

* github > settings > developer settings > personal access tokens > access tokens (classic) > generate new token (classic)
  * `write:packages`
* copy the token
* set as `GHCR_TOKEN` in actions secrets on this repo