# busgopher

![logo](./docs/logo.png)

## A simple terminal client for sending messages to Azure Service Bus.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)

The tool currently has one responsibility: sending messages to [Azure Service Bus](https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-messaging-overview).

## Usage

Busgopher may be run in two modes: CLI and GUI.

### CLI

The CLi mode is dedicated to run fast the tool without UI. To use CLI mode run busgopher with arguments that will specify: connection, destination and message. Connection, destination and message are selected by providing their name.

```sh
./busgopher --msg="test-message" --conn="demo" --dest="test-queue"
```

Don't worry if you miss argument or pass a wrong name. Busgopher will provide an error message. 

```sh
./busgopher --msg="unknown" --conn=dev --dest=test-queue
Started headless mode with connection: demo, destination: test-queue, message: unknown
[2024-10-08 20:04:37]: [Info] Selected connection: dev (demo.servicebus.windows.net)
[2024-10-08 20:04:37]: [Info] Selected destination: test-queue
[2024-10-08 20:04:37]: [Error] Can't find message with name: unknown
[2024-10-08 20:04:37]: [Error] Message not selected!
```

### GUI

## Configuration

The busgopher configuration is divided into two parts: connections & messages. At the moment to edit/set the configuration you have to manually edit config.json file located next to busgopher executable.

### Connections

The first one is responsible for configuring the connection to Azure Service Bus (ASB). The connection uses [DefaultAzureCredentials](https://learn.microsoft.com/en-gb/dotnet/azure/sdk/authentication/credential-chains?tabs=dac#usage-guidance-for-defaultazurecredential) feature to authorize access to ASB.

Available properties:
- Name - user friendly connection name. Will be used to select connection;
- Namespace - Azure Service Bus namespace;
- Destinations - list of available entitites (both queues and topic) which may be select to send a message to;

Sample connection section:
```json
[
    {
        "Name": "dev",
        "Namespace": "<<namespace>>.servicebus.windows.net",
        "Destinations": [
            "test-queue",
            "test-topic"
        ]
    }
]
```

### Messages

The second part of the configuration - messages, define messages that will be send to ASB. We may define both built-in and custom message properties. More about properties in the Features section.

Besides properties user may define message body. Message body has no structure limitations. The body length is limited by your ASB configuration in Azure. You can benefit from message template engine while defining message's body. More about message body template engine in the Features section.

Sample message section:

```json
[
    {
        "name": "message-name",
        "body": "Sample body with template engine generated at {{utcNow}}",
        "subject": "subject",
        "customProperties": { 
            "another": "false" 
        }
    }
]
```

### Sample config

```json
{
    "Connections": [
        {
            "name": "dev",
            "namespace": "demo.servicebus.windows.net",
            "destinations": [
                "test-queue",
                "test-topic"
            ]
        }
    ],
    "Messages": [
        {
            "name": "another test",
            "body": "another test {{utcNow}}, {{generateUUID}}",
            "subject": "json subject",
            "customProperties": { 
                "another": "false" 
            }
        }
    ]
}
```

## Features

### Message body templates engine

Bustopher provides a simple template engine that enables message body generation according to a defined template. The built-in template engine is based on the Golang [text/templates package](https://pkg.go.dev/text/template). Bustopher provides a set of predefined functions. 

To use the engine, embed a predefined function in the saved message body, like:

```json
    {
        "name": "Engine presentation",
        "body": "Random UUID: {{generateUUID}} generated at {{utcNow}} "
    },
```

#### Predefined functions

- utcNow
Gets current UTC now date time from machine and returns it in RFC333 format. Usage:
```
Message created at {{utcNow}}.

Message created at 2024-10-06T19:34:39Z.
```

- generateUUID
Generates random UUID. Usage:
```
This is random UUID: {{generateUUID}}.

This is random UUID: 69a17b86-68d7-4e59-bb2f-09b3590135c8.
```

### Message properties

Busgopher supports defining messages built in and custom properties that consumers may use. Supported built in properies are:
- CorrelationID
- MessageID
- ReplayTo
- Subject

To define messages' properties just define them in the messages.json file like:

```json
    {
        "name": "Properties sample",
        "body": "Focus on properties please",
        "subject": "Custom subject",
        "replyTo": "def-null"
        "customProperties": { 
            "isCustom": "true" 
        }
    }
```

## Your Feedback

Add your issue here on GitHub. Feel free to get in touch if you have any questions.

## Code of Conduct

This project has adopted the code of conduct defined by the [Contributor Covenant](https://www.contributor-covenant.org/) to clarify expected behavior in our community.

