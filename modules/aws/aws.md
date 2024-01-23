import { Callout } from 'nextra/components';

# aws

<Callout type="info" emoji="ℹ️">
  This module requires that Risor has been compiled with the `aws` Go build tag.
  When compiling **manually**, [make sure you specify `-tags aws`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source).
</Callout>

The `aws` module exposes a simple interface that wraps the AWS SDK v2 for Go.
Create a client by providing the name of the service you want to use. All API
calls for that service are then made available on the client.

Speciying a config is optional, but can be used to configure the client.

## Functions

### client

```go filename="Function signature"
client(service string, config map | aws.config) aws.client
```

Creates a new AWS client for the given service. The `config` parameter is
optional and can be used to configure the client. All service API calls are
made available on the client.

```go copy filename="Example"
>>> aws.client("s3")
aws.client(service=s3, region=us-east-1)
>>> aws.client("ec2")
aws.client(service=ec2, region=us-east-1)
>>> aws.client("ec2", {region: "us-west-2"})
aws.client(service=ec2, region=us-west-2)
```

### config

```go filename="Function signature"
config(config map) aws.config
```

Creates a new AWS config with the given configuration map.

```go copy filename="Example"
>>> aws.config({region: "us-east-2"})
aws.config(region=us-east-2)
```

Available configuration options:

```json
{
  "region": "us-east-1",
  "credentials": {
    "key": "AKID",
    "secret": "SECRET",
    "session": "SESSION_TOKEN"
  },
  "profile": "custom_profile",
  "credentials_files": ["test/credentials"],
  "config_files": ["test/config"]
}
```

## S3 Client Usage

```go copy
>>> s3 := aws.client("s3")
>>> s3.list_buckets()["Buckets"]
[{"CreationDate "2023-07-26T01:16:19Z", "Name "example-12345"}]
>>> s3.create_bucket({Bucket: 'test-{rand.int()}'})
{"Location": "/test-2769212968479898940", "ResultMetadata": {}}
```

## Services

Support for these services is built into the Risor CLI:

- `apigatewayv2`
- `athena`
- `backup`
- `cloudformation`
- `cloudfront`
- `cloudtrail`
- `cloudwatch`
- `cloudwatchlogs`
- `ddb`
- `ebs`
- `ec2`
- `ecr`
- `ecs`
- `eks`
- `elasticache`
- `elasticsearchse`
- `eventbridge`
- `firehose`
- `glue`
- `iam`
- `kinesis`
- `kms`
- `lambda`
- `ram`
- `rds`
- `redshift`
- `route53`
- `s3`
- `secretsmanager`
- `sesv2`
- `sfn`
- `sns`
- `sqs`
- `sts`
- `wafv2`
- `xray`
