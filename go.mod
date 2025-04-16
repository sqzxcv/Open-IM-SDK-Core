module github.com/openimsdk/openim-sdk-core/v3

go 1.23.0

toolchain go1.23.3

require (
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.5.3
	github.com/jinzhu/copier v0.4.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.33.0 // indirect
	gorm.io/driver/sqlite v1.5.7

)

require golang.org/x/net v0.37.0

require (
	github.com/coder/websocket v1.8.12
	github.com/google/go-cmp v0.7.0
	github.com/openimsdk/protocol v0.0.72
	github.com/openimsdk/tools v0.0.24
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/sqzxcv/glog v0.0.0-20240903100204-ee7f108f80dd
	golang.org/x/image v0.25.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.62.1 // indirect
)

//replace github.com/coder/websocket => /Users/shengqiang/Documents/Codes/wasm-websocket
//replace github.com/OpenIMSDK/protocol => github.com/openimsdk/protocol v0.0.72

replace github.com/openimsdk/protocol => /Users/shengqiang/Documents/Codes/protocol

replace github.com/openimsdk/tools => github.com/sqzxcv/openim-tools v0.0.0-20241214123830-a19f8c9140d2
