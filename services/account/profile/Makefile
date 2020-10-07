.PHONY: mocks
mocks:
	# mockgen -source=token.go -package=profile -destination=token_mock.go
	mockgen -source=profilestore.go -package=profile -destination=profilestore_mock.go
	mockgen -source=service.go -package=profile -destination=service_mock.go
	mockgen -source=client.go -package=profile -destination=client_mock.go
