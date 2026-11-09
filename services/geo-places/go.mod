module github.com/makskozlov/atlasflow/services/geo-places

go 1.24

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/makskozlov/atlasflow/pkg v0.0.0
	github.com/makskozlov/atlasflow/proto v0.0.0
	github.com/stretchr/testify v1.10.0
)

replace github.com/makskozlov/atlasflow/pkg => ../../pkg

replace github.com/makskozlov/atlasflow/proto => ../../proto