# gcloud-pubsub-emulator-helper

> [!WARNING]
> This project was created as a side project to learn Go and to solve a specific problem encountered at work. Therefore, it should be treated as a **toy project**.

## Overview
This project serves as a helper for the official GCloud Pub/Sub emulator, providing additional functionality while keeping the core dependency-free.

- [GCloud Emulator for Pub/Sub](https://cloud.google.com/pubsub/docs/emulator)

### Future Plans
In the future, I plan to add a Web User Interface to visualize basic data from the emulator. This would act as an additional entry point and may introduce some dependencies.

## Features

> [!NOTE]
> The following features are designed based on the Pub/Sub REST API documentation. Some of them might be removed if they are not supported by the emulator.

- [X] Load JSON configuration file
- [X] Basic sync between the emulator and the provided configuration
- [X] Support for Labels in Topics
- [X] Support for Labels in Subscriptions
- [X] Support for Message Storage Policy
- [X] Support for KMS Key Name
- [ ] Support for Schemas
- [X] Support for State Response (Emulator returns a dumb empty value)
- [ ] Additional Web GUI build entry

🔗 [GCloud Pub/Sub REST API Documentation](https://cloud.google.com/pubsub/docs/reference/rest)

## Building the Executable
To generate the loader executable, run:

```bash
make build
```

This command will create an executable called `./basicLoader`, which is the core component of this project.

## Running the Executable

### Configuration Setup
Create a JSON configuration file. The default file name is `config.json`, but you can specify a different file if needed.

### Executable Arguments
The following command-line arguments are available:

- **`-config`** *(string, default: `./config.json`)* - Path to the JSON configuration file.
- **`-host`** *(string, optional)* - Overrides the host specified in the configuration file.
- **`-help`** *(boolean, default: `false`)* - Displays the help message and exits.

#### Example Usage
```sh
# Run with a specific configuration file
./basicLoader -config=/path/to/config.json

# Run with a custom host
./basicLoader -host=127.0.0.1:8085

# Show help message
./basicLoader -help
```

#### Notes
- If `-help` is provided, the application prints the available options and exits.
- If no `-config` argument is provided, the application defaults to `./config.json`.
- If an invalid `-host` is provided, the application exits with an error.

## Configuration File

### JSON Structure
```json
{
  "delayBeforeStartupCheckMs": 0,
  "avoidStartupCheck": false,
  "startTimeoutMs": 30000,
  "timeBetweenStartupChecksMs": 200,
  "projects": [
    {
      "name": "project-name",
      "topics": [
        {
          "name": "topic.name.1",
          "subscriptions": [
            {
              "name": "subscription.1.for.topic.name.1"
            },
            {
              "name": "subscription.2.for.topic.name.1"
            }
          ]
        }
      ]
    }
  ]
}
```

### Configuration Fields

#### Global Settings
- **`delayBeforeStartupCheckMs`** *(integer)* - Delay in milliseconds before the startup check.
- **`avoidStartupCheck`** *(boolean)* - If `true`, skips the startup check.
- **`startTimeoutMs`** *(integer)* - Maximum wait time (in milliseconds) for the emulator to start.
- **`timeBetweenStartupChecksMs`** *(integer)* - Time interval (in milliseconds) between startup checks.

#### Project Settings
The `projects` array defines the Pub/Sub projects.

- **`name`** *(string)* - Name of the project.
- **`topics`** *(array)* - List of topics within the project.
  - **`name`** *(string)* - Name of the topic.
  - **`labels`** *(map[string]string, optional)* - Labels added to the topic.
  - **`messageStoragePolicy`** *(MessageStoragePolicy, optional)* - Policy that should be applied for message storage.
    - **`allowedPersistenceRegions`** *([]string, optional)* - [Google Cloud Region's IDs](https://cloud.google.com/about/locations).
    - **`enforceInTransit`** *(bool, optional)* - If `true`, `allowedPersistenceRegions` is also used to enforce in-transit guarantees for messages.
  - **`kmsKeyName`** *(string, optional)* - The resource name of the Cloud KMS CryptoKey to be used to protect access to messages published on this topic.
  - **`messageRetentionDuration`** *(string, optional)* - AVOID. This field does not seem to be accepted by the emulator but it exists in the REST API.
  - **`subscriptions`** *(array, optional)* - List of subscriptions for the topic.
    - **`name`** *(string)* - Name of the subscription.
    - **`labels`** *(map[string]string, optional)* - Labels added to the subscription.
  - **`ingestionDataSourceSettings`** *(IngestionDataSourceSettings, optional)* - Configuration for external ingestion sources.
    - **`platformLogsSettings`** *(PlatformLogsSettings, optional)* - Configuration for platform log ingestion.
      - **`severity`** *(string)* - The severity level of logs to ingest (e.g., `INFO`, `WARNING`, `ERROR`).
    - **`awsKinesis`** *(AwsKinesisSettings, optional)* - Configuration for AWS Kinesis ingestion.
      - **`state`** *(string)* - The state of the AWS Kinesis ingestion (e.g., `ACTIVE`, `KINESIS_PERMISSION_DENIED`).
      - **`streamArn`** *(string)* - The ARN of the AWS Kinesis stream.
      - **`consumerArn`** *(string)* - The ARN of the AWS Kinesis consumer.
      - **`awsRoleArn`** *(string)* - The ARN of the AWS IAM role used for access.
      - **`gcpServiceAccount`** *(string)* - The Google Cloud service account associated with ingestion.
    - **`cloudStorage`** *(CloudStorageSettings, optional)* - Configuration for Cloud Storage ingestion.
      - **`state`** *(string)* - The state of the Cloud Storage ingestion (e.g., `ACTIVE`, `BUCKET_NOT_FOUND`).
      - **`bucket`** *(string)* - The Cloud Storage bucket used for ingestion.
      - **`minimumObjectCreateTime`** *(string)* - The minimum object creation time to consider.
      - **`matchGlob`** *(string, optional)* - A glob pattern to match object names.
      - **`textFormat`** *(TextFormatSettings, optional)* - Configuration for text format ingestion.
        - **`delimiter`** *(string)* - The delimiter used in text format files.
      - **`avroFormat`** *(object, optional)* - If present, indicates Avro binary format ingestion.
      - **`pubsubAvroFormat`** *(object, optional)* - If present, indicates ingestion of Pub/Sub exported Avro messages.

## Internal Functionality

### 1️⃣ Loading the Configuration
- `LoadConfigurationFromFile(filepath string) (Configuration, error)`
  - Reads the JSON configuration file and unmarshals it into the `Configuration` struct.
  - Applies default values if necessary.
  - Exits the application if an invalid host is provided.

### 2️⃣ Syncing with the Emulator
- `Sync(client utils.Client)`
  - Ensures the emulator reflects the provided configuration.
  - Waits for the emulator to be available if `avoidStartupCheck` is `false`.
  - Deletes existing topics and subscriptions before applying the new configuration.
  - Creates new topics and subscriptions based on the configuration.

## Working with this repository
We use `pre-commit` in order to have all the files checked out and testing
passed before commiting.

To use it, after clonning the repo and [installing pre-commit](https://pre-commit.com/#install), you need to execute:
```bash
pre-commit install
```

By default it will always run the commands in the files you are commiting, but
you could execute it in every file with the following command:
```bash
pre-commit run --all-files
```

For more information: [pre-commit](https://pre-commit.com/)

## Legal Information

This project is not affiliated with, endorsed by, or officially connected to Google or Google Cloud Platform.

Google Cloud Pub/Sub, the Google Cloud Emulator, and related trademarks, service marks, and logos are the property of Google LLC.

For more information on Google Cloud Pub/Sub, visit the official documentation: [Google Cloud Pub/Sub](https://cloud.google.com/pubsub).
