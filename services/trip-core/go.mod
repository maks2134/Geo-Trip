module github.com/makskozlov/atlasflow/services/trip-core

go 1.25.4

require (
	connectrpc.com/connect v1.18.1
	github.com/google/uuid v1.6.0
	github.com/makskozlov/atlasflow/pkg v0.0.0
	github.com/makskozlov/atlasflow/proto v0.0.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/net v0.34.0
	google.golang.org/protobuf v1.36.5
)

replace (
	github.com/makskozlov/atlasflow/pkg => ../../pkg
	github.com/makskozlov/atlasflow/proto => ../../proto
)