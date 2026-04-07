# WHITE.archive

Encrypted incremental backup utility files to S3-compatible storage.

## Usage

```bash
white-archive -mode=backup -dir=/data
white-archive -mode=restore -dir=/data
```

| Flag | Default | Description |
|------|---------|-------------|
| `-mode` | — | `backup` or `restore` |
| `-dir` | `/data` | Directory to backup / restore into |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `CRYPTO_KEY` | Encryption key |
| `S3_ENDPOINT` | S3-compatible endpoint (e.g. `s3.eu-central-1.backblazeb2.com`) |
| `S3_ACCESS_KEY` | S3 access key |
| `S3_SECRET_KEY` | S3 secret key |
| `S3_BUCKET` | Bucket name |