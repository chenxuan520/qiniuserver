---
name: qiniuserver
description: Use this skill whenever an agent needs to upload, list, preview, or delete files on a qiniuserver instance (the self-hosted 七牛云 Kodo file bed in this repository). Trigger it for requests like "上传到图床", "把文件传到七牛", "list the files on qiniuserver", "delete this uploaded file", or any task that calls qiniuserver's /api/upload, /api/upload_url, /api/list, /api/delete, or /api/info endpoints. Also use it when you are working inside this repo and need the exact request shapes, auth rules, and response envelopes instead of guessing them from the frontend code.
---

# qiniuserver

qiniuserver is a single Go binary that serves a small web UI plus an HTTP API which proxies file operations to 七牛云 Kodo. This skill covers calling that API as a client. For changing the Go code itself, follow `AGENTS.md` instead.

## 1. Find the base URL

Use the origin of the web UI, for example `http://127.0.0.1:5200`. If the user does not provide one:

1. Read `config/config.json` and combine `host.ip` and `host.port`.
2. Or read the URL the user is already browsing.

All endpoints live under `/api`. Never hardcode a Qiniu domain; ask `/api/info` for it.

## 2. Authentication

Set the `Authorization` header to the raw value of `host.password`:

```
Authorization: <password>
```

Important details:

- Do **not** add a `Bearer ` prefix. The middleware compares the header value to the configured password directly (`src/middlerware/auth.go:33-36`).
- If `host.password` is empty in the config, authentication is disabled and no header is needed.
- `401` means the password is wrong. `403` means attempts are being throttled for 3 seconds after a failure. Do not retry in a tight loop.

## 3. Response envelope

All handlers return `{ "code": number, "data": any }` (`src/controller/response/res.go`).

- Success: HTTP 200, `code: 200`, `data` is the payload below.
- Business error: HTTP 400 with `code: 400` and `data` containing a human-readable string.
- Auth errors: HTTP 401/403 with `{ "code": ..., "msg": "..." }`, no `data`.

## 4. Endpoints

### GET /api/info

Returns the configured upload prefix and CDN domain.

```bash
curl -s "$BASE/api/info" -H "Authorization: $PASSWORD"
# {"code":200,"data":{"path":"upload_path","domain":"https://cdn.example.com"}}
```

`path` is the Kodo key prefix used for uploads and listing. `domain` is used to build the final file URL.

### POST /api/upload

Multipart upload. The field name is exactly `file`.

```bash
curl -s "$BASE/api/upload" \
  -H "Authorization: $PASSWORD" \
  -F "file=@/absolute/path/to/photo.png" \
  -F "file=@/absolute/path/to/doc.pdf"
# {"code":200,"data":{"url":"https://cdn.example.com/upload_path/photo.png"}}
```

Behavior and constraints:

- The client filename may be rewritten by the server's `qiniu.file_name` template (`%d` timestamp, `%f` original name, `%r` random string), so read the returned `url` instead of assuming the key.
- The server first writes the upload to its working directory using the client-provided filename. Send a plain basename, never a path such as `../../etc/x`.
- If the returned `url` has no scheme because `qiniu.domain` lacks one, prepend `https://` before using it as a link.

### POST /api/upload_url

Downloads a URL server-side and stores the content.

```bash
curl -s "$BASE/api/upload_url" \
  -H "Authorization: $PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{"path":"https://example.com/image.png"}'
# {"code":200,"data":{"url":"https://cdn.example.com/upload_path/xxxxxx.url"}}
```

Only use this for URLs the user trusts and controls. The server fetches whatever URL it is given, so an arbitrary internal address is reachable from the server.

### POST /api/list

Paginated listing. All three body fields are optional except `path`.

```bash
curl -s "$BASE/api/list" \
  -H "Authorization: $PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{"path":"upload_path","marker":"","limit":50}'
# {"code":200,"data":{"list":[{"key":"photo.png","fsize":1234,"putTime":...}],"next_marker":"...","has_next":true}}
```

- `path` is the prefix, normally the `path` returned by `/api/info`.
- Returned `key` values are relative to that prefix.
- `limit` is clamped to the 1..1000 range, defaulting to 50 when outside it.
- Loop until `has_next` is false, passing the previous `next_marker` as `marker`.

```bash
marker=""
while :; do
  resp=$(curl -s "$BASE/api/list" -H "Authorization: $PASSWORD" \
    -H "Content-Type: application/json" \
    -d "{\"path\":\"$PREFIX\",\"marker\":\"$marker\",\"limit\":1000}")
  echo "$resp" | jq -r '.data.list[].key'
  marker=$(echo "$resp" | jq -r '.data.next_marker // empty')
  [ "$(echo "$resp" | jq -r '.data.has_next')" = "true" ] || break
done
```

### POST /api/delete

Deletes one object. `path` must be the full Kodo key, including the upload prefix.

```bash
curl -s "$BASE/api/delete" \
  -H "Authorization: $PASSWORD" \
  -H "Content-Type: application/json" \
  -d '{"path":"upload_path/photo.png"}'
# {"code":200,"data":null}
```

Because `/api/list` returns keys relative to the prefix, join `PREFIX + "/" + key` before deleting. Deletion is irreversible; confirm with the user before deleting anything that was not just uploaded in this task.

### POST /api/upload_test

This endpoint does not upload; it only computes and returns a URL. Never use it to store files.

## 5. Recommended workflow

1. `GET /api/info` to obtain `path` and `domain`. Fail fast if the password is wrong.
2. Upload with `POST /api/upload`, or `POST /api/upload_url` when the source is a URL.
3. Report the returned `url` to the user. Treat it as the source of truth.
4. For listing, page through `POST /api/list` until `has_next` is false.
5. For deletion, reconstruct the full key and confirm destructive intent first.

## 6. Failure handling

- `data` is a string on business errors; surface it rather than returning just the status code.
- A `401` means credentials, not a transient failure. Do not blindly retry.
- `403` is a global 3-second throttle, so wait at least 3 seconds before the next attempt.
- A network failure on `/api/upload_url` is about the remote URL, not the qiniuserver instance.
