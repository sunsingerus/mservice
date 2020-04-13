module github.com/sunsingerus/mservice

go 1.13

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.5
	github.com/google/go-cmp v0.2.0
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.28.1
	gopkg.in/square/go-jose.v2 v2.4.1 // indirect
)

replace golang.org/x/oauth2 => github.com/sunsingerus/oauth2 v0.0.0-20200410181841-d7afaacd4cbe
