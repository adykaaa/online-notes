## Online NoteZ!

A full-stack application that will be built up from scratch using Go for the backend, and React for the frontend. It's still a WIP.

Plans:
- auth handled via PASETOs in HTTP only cookies
- frontend to be done with ChakraUI
- app will be a monolith (with 1 Postgre DB)
- logging with Zerolog cuz it's fast and supports JSON
- HTTP backend will use Chi as it's really close to the stdlib, and route grouping honestly rocks
- I'll implement unit testing as I go along (I'm working in a way that's similar to TDD, but sometimes I get lazy)
- Integration testing will deffo be there
- Caching with Redis
- Notes to PDF feature -> will be done with RabbitMQ for adding some event-driven madness
- OpenTelemetry will be implemented for sure.
- I'll forward logs to Elastisearch just for fun
- the whole app's CI will be done with GHA, and I'll deploy it to K8s

REALLY REALLY future plan: make it cloud native, mostly using Lambda functions, RDS, and an API gateway