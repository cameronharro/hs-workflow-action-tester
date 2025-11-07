# Goal

- Enable a faster development loop for HubSpot Custom Workflow Actions
- Create a CLI tool that integrates with HubSpot project files
- Make automated testing possible for workflow action CI

## Features

- Parse HubSpot workflow action definition files and configurable profiles
- Accept data test cases in CSV format
- Impersonate HS server sending option and execution requests to app back end
- Compare responses and callbacks to expected values and log/exit accordingly

## Configuration/Flags

- Environment file (for App secrets, etc.)
- Definition/call to test
- HS profile (local dev/staging/prod) to use to fill out definition
- Test case file(s?)
