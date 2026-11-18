# WHITE.archive

Encrypted incremental backup utility to S3-compatible storage.

All parameters required. 
  Parameters can be specified either through environment variables or through flags.

## Environment Variables

| Flag             | Variable        | Description                        |
|------------------|-----------------|------------------------------------|
| `-mode`          | `MODE`          | Run mode, backup or restore        |
| `-directory`     | `DIRECTORY`     | Direcrory for backup or restore to |
| `-crypto-key`    | `CRYPTO_KEY`    | Encryption key                     |
| `-s3-endpoint`   | `S3_ENDPOINT`   | S3-compatible endpoint             |
| `-s3-access-key` | `S3_ACCESS_KEY` | S3 access key                      |
| `-s3-secret-key` | `S3_SECRET_KEY` | S3 secret key                      |
| `-s3-bucket`     | `S3_BUCKET`     | Bucket name                        |


> useless shit
