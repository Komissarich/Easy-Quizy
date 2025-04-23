# Easy Quizy

---
## Table Of Contents

1. [About](#about)
2. [Architechture](#architechture)
3. [Stack](#stack)
4. [Microservices APIs](#microservices-apis)
   1. [Authentification service](#authentification-service)
   2. [Quiz service](#quiz-service)
   3. [Statistics service](#statistics-service)

---
## About

This is a social network for quiz lovers. Here you can take quizzes from other users or create your own, find friends, compete in the overall rating of authors or players.

## Architechture 

scheme.png

## Stack

- Backend: Golang, Protobuf  
- Frontend: Vue.js
- DB: PostgreSQL, Redis
- DevOps: Docker

## Microservices APIs

!!!include(auth_service/README.md)!!!

!!!include(quiz_service/README.md)!!!

!!!include(stat_service/README.md)!!!