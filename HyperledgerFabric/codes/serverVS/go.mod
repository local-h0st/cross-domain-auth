module serverVS

go 1.20

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/hyperledger/fabric-gateway v1.2.2 // indirect
	github.com/hyperledger/fabric-protos-go-apiv2 v0.3.0 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require myrsa v0.0.0-00010101000000-000000000000
replace myrsa => ../toolPackages/myrsa
require msgs v0.0.0-00010101000000-000000000000
replace msgs => ../toolPackages/msgs