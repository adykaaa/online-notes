# Online NoteZ!

This project aims to be the most over-engineered fullstack note-taking application, using React for frontend and Go for backend.

![](https://github.com/adykaaa/online-notes/blob/main/web/public/onlinenotes.gif)

## Over-engineered, how? you ask
The domain logic is pretty simple. A user can create notes, update them, and delete them, nothing to ride home about. The aim of this whole project is to showcase (and learn) how to:
- do unit testing, and integration testing in Go (TBD)
- use Chi as a router
- use SQLC for database operations
- gomock for generating cool mocks
- PASETO user authentication using secure and HTTP-only Cookies
- use docker-compose to stand the whole dev environment up
- use custom middlewares (for logging, authentication, etc.)
- use Viper for configuration management through environment variables
- use Zerolog for logging, that's being passed around in the request context
and many more. 

I am NOT stating these are the most idiomatic ways of doing things, and everything in this repo is subject to change.

## Plans

So basically after the "base" is done, where a user can register, login, create/delete/update notes - with all unit tests, I'm planning to:
- add a "Chat" feature using websockets and Redis
- add backend caching using Redis
- use opentracing and opentelemetry
- deploy everything into a local K8s cluster
- stand up Elastisearch and Kibana for searching through the logs
- migrate everything to AWS serverless, and document the whole journey

