# SAM deploy example:
sam build
sam deploy --config-file samconfig.toml --region ap-southeast-1 --tags "auto-delete=no" --resolve-s3 