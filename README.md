# HubSpot Workflow Tester

## Overview

Documentation for HubSpot's Custom Workflow Actions can be found [here](https://developers.hubspot.com/docs/api-reference/latest/automation/workflow-actions/custom-action-reference)

Testing HubSpot's custom workflow actions is cumbersome. Using this project, integration testing between custom workflow action definitions and your application can be done locally.

Upon making changes to either your application or your HubSpot custom action definition, they should be tested to ensure that the requests, responses, and callbacks match each others' expectations.

### Existing Process
1. Create additional "test" HubSpot project or action definition so that regressions do not affect existing users
2. Create test accounts and authenticate them with your app
3. Create workflows using your custom action and records matching your app's requirements
4. Manually run records through the workflows, then verify correctness by hand
5. Repeat steps 3 and 4 for each change you make, and also make sure there are no regressions
6. Deploy changes to "real" HubSpot project

This process is slow and painful, especially when manually running workflows and verifying their results.

### Updated Process
1. ~~Create additional "test" HubSpot project or action definition~~ Local testing allows you to skip the "test" project because you don't need to deploy a HubSpot action configuration
2. ~~Create test accounts and authenticate them with your app~~ Only necessary if your app depends on actual HubSpot accounts (e.g. fetching data from the APIs for the sending portal)
3. ~~Create workflows using your custom action and records matching your app's requirements~~
4. ~~Manually run records through the workflows, then verify correctness by hand~~ specify test cases and expected results up front, run and verify locally in a single command
5. ~~Repeat steps 3 and 4 for each change you make, and also make sure there are no regressions~~ Incorporate test script into CI, never forget a test case
6. Deploy changes to HubSpot project

## Usage

The tester is a CLI tool that takes a file of json-formatted test cases, a HubSpot project that contains custom workflow action definitions, and a few configuration parameters. It:
- Sends the parsed test cases to the specified address of your app server (probably running locally, but not necessarily)
- Listens for responses or async callbacks
- Evaluates them against the test expectations

### Features

#### Authentication

The tester signs the HTTP request with x-hubspot-signature headers for both v2 and v3 (configurable by injecting value for client secret)

#### JS functions [beta]

Supports user-provided Javascript snippets to mutate request/response values to adapt between your application and HubSpot's systems

This feature passes internal tests, but the HubSpot's documentation is sparse and I don't use this feature myself. Contributions/issues welcome

#### Async callbacks

The tester sets up an HTTP server to receive asynchronous callbacks from your application (configurable by injecting listening Port number)

Also requires making your application's callback domain configurable so that callbacks can be pointed to this server instead of https://api.hubapi.com

#### Evaluation

Supports evaluation against HubSpot's Execution Rules, which drive messages displayed to the end user

### Quickstart

#### Dependencies

- Go (version 1.26 minimum)
- Deno (version 2.5.6 minimum) -- for sandboxed Javascript functions

#### Installation

Clone this repo locally: `git clone https://github.com/cameronharro/hs-workflow-action-tester`

Install Go dependencies: `go mod download`

Install CLI globally: `go install .` (requires Go binaries to be accessible in your PATH)

#### Test

Run Go tests: `go test ./...`

Run the dummy server in one process: `node server.js`

Run the tester (example config): `go run . -cases testCases.json -clientSecret 12345`

### Configuration

Once you have the CLI installed globally, you can invoke the tool as `hs-workflow-tester [FLAGS]`

1. Navigate to your HubSpot project
2. Create some test cases as shown below
3. Run your own app locally
4. Invoke the tester pointing at your local app

#### Tests

Tests are stored in a .json file as an array of tests structured as the below. See ./testCases.json for an example

```json
{
    "testLabel": "1", // string, required: the tester's label, is not exposed
    "actionUID": "test_action", // string, required: the HubSpot project action definition's ID
    "portalID": 111, // int, required: The HubSpot portal that is being simulated
    "type": "action", // enum ("action", "option"), required: The type of behavior being tested
    "test": {} // (ActionTest, OptionTest), required: The kind of data specified by the type field above
}
```

#### ActionTest

```json
{
    "actionURL": "http://localhost:3000", // string: the url to send the action request to. Defaults to configured value from the action definition
    "inputFields": {
        "label": "Brian Halligan",
        "value": "1234"
    }, // object: no nested fields allowed, but otherwise any primitive JSON type may be paired with a string key
    "objectID": 1234, // int, required: The ID of the record being simulated
    "objectType": "CONTACT", // string, required: The type of the record being simulated
    "expectedExecutionRule": "success" // string: The expected execution rule. Whether the response + action definition yield this rule determines whether the test succeeds or not.
}
```

#### OptionTest

```json
{
    "optionsURL": "http://localhost:3000", // string: the url to send the action request to. Defaults to configured value from the action definition
    "objectTypeId": "0-1", // string, required: The numeric ID of the object type that the simulated workflow runs
    "inputFieldName": "options_field", // string, required: the name of the field for which to fetch options
    "inputFields": {
        "label": {
            "type": "STATIC_VALUE",
            "value": "Brian Halligan"
        }
    }, // object: each known input field needs an object with "type" and either "propertyName" or "value" values
    "expectedOptions": [
        {
          "label": "Test",
          "value": "test"
        }
    ] // array: The expected options; order matters.
}
```

#### CLI Flags

- cases (required): relative path to the test case file
- clientSecret (required): string to be used when signing requests as HubSpot does
- port: number of port on which the async callback server should listen
- project: relative path to the HubSpot containing the action definition to be tested
- timeout: duration string as described in [the Go docs](https://pkg.go.dev/time#ParseDuration)
