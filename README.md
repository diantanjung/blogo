## About blogo
An simple personal blog application for author to write articles. This is example implementation of microservice with Clean Architecture in Go (Golang) project.

### This project has 4 Domain layer :
- Model/Entity Layer
  model is a struct reflecting our data object from / to database. Models should only define data structs, no other functionalities should be included here.
- Repository Layer
  Repository Layer is where the implementation of Model layer. All queries and data operation from / to database should happen here.
- Usecase/Service Layer
  Usecase Layer is where the business logic lies on, it handles delivery layer request and fetch data from Repository layer it needs and run their logic to satisfy what delivery layer expect the service to return.
- Delivery/Controller Layer
  Delivery Layer are the handler of all requests coming in, to the router.


### List service
- [x] User service
- [ ] Article service

### Tools and library
- All libraries listed in [`go.mod`]
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.

### Reference :
- https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- https://github.com/gsvaldevieso/go-dream-architecture
- https://github.com/caohoangnam/go-clean-architecture
- https://github.com/bxcodec/go-clean-arch
- https://github.com/irahardianto/service-pattern-go