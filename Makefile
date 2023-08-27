generate-mocks:
		mockgen -package=mocks -source=internal/domain/service/api.go -destination=mocks/services_mock.go
		mockgen -package=mocks -source=internal/repository/api.go -destination=mocks/repositories_mock.go

