.PHONY: protos

protos:
	protoc protos/currency.proto --go_out=. --go-grpc_out=.

list-services:
	grpcurl --plaintext localhost:9092 list

list-currency-methods:
	grpcurl --plaintext localhost:9092 list Currency        

rate-request-details:
	grpcurl --plaintext localhost:9092 describe .RateRequest

example-execute-rate-request:
	grpcurl --plaintext -d '{"base": "GBP", "destination": "USD"}' localhost:9092 Currency/GetRate
