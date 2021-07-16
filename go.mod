module github.com/certusone/janus

go 1.16

require (
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/tendermint/tendermint v0.34.0-rc3.0.20200907055413-3359e0bf2f84
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200921161331-205a656cc58b
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20200806022845-90696ccdc692 // indirect
	google.golang.org/genproto v0.0.0-20200806141610-86f49bd18e98 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace go.etcd.io/etcd/v3 => github.com/hendrikhofstadt/etcd/v3 v3.3.0-rc.0.0.20200923163016-ac81520e9e28
