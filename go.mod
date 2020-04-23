module github.com/sunsingerus/mservice

go 1.13

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.5
	github.com/google/uuid v1.1.1
	github.com/maxatome/go-testdeep v1.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1 // indirect
	google.golang.org/grpc v1.28.1
)

replace golang.org/x/oauth2 => github.com/sunsingerus/oauth2 v0.0.0-20200410181841-d7afaacd4cbe
