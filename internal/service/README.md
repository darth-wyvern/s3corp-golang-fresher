# internal/service

Business logic layer:

- Define Service interface: an interface declare all function of the entity service.
- Define Service implement: a struct which implement the `service interface`,must be dependent on `repository interface`, and call any function of it.
- Define Service mock: a struct which implement the `service interface`, but it is fake to test the handler layer.

- Control business logic which contain many step
- Responsible for checking business rules, logical constraints, and performing tasks.
