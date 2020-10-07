.PHONY: mocks
mocks:
	mockgen -source=token.go -package=account -destination=token_mock.go
	mockgen -source=accountstore.go -package=account -destination=accountstore_mock.go
	mockgen -source=service.go -package=account -destination=service_mock.go
	mockgen -source=client.go -package=account -destination=client_mock.go
