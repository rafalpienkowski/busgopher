# busgopher

![logo](./docs/logo.png)

## A simple terminal client for sending messages to Azure Service Bus.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)

The tool currently has one responsibility: sending messages to (Azure Service Bus)[https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-messaging-overview].

## Usage

### UI vs CLI mode

## Configuration

## Features

### Message body templates engine

Bustopher provides a simple template engine that enables message body generation according to a defined template. The built-in template engine is based on the Golang (text/templates package)[https://pkg.go.dev/text/template]. Bustopher provides a set of predefined functions. 

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

This project has adopted the code of conduct defined by the (Contributor Covenant)[https://www.contributor-covenant.org/] to clarify expected behavior in our community.

