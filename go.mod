module github.com/atomix/kubernetes-simulations

go 1.12

require (
	github.com/atomix/go-client v0.0.0-20200207221255-96f6ea5d353d
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/onosproject/onos-test v0.0.0-20200212220529-60ed5cef794f
)

replace github.com/atomix/go-client => ../atomix-go-client

replace github.com/onosproject/onos-test => ../onos-test
