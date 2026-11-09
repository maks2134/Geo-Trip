module github.com/makskozlov/atlasflow/services/gateway-auth

go 1.25.4

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/makskozlov/atlasflow/pkg v0.0.0
	github.com/makskozlov/atlasflow/proto v0.0.0
	golang.org/x/crypto v0.35.0
)

replace github.com/makskozlov/atlasflow/pkg => ../../pkg
replace github.com/makskozlov/atlasflow/proto => ../../proto