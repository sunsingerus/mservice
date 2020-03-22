// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service_auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"strings"

	log "github.com/golang/glog"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "No metadata provided")
	errMissingToken    = status.Errorf(codes.Unauthenticated, "No authorization token provided")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "Invalid token")

	jwtPublicKey []byte
)

// In case of failed authorization, the interceptor blocks execution of the handler and returns an error.
// type grpc.StreamServerInterceptor
func streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	log.Infof("streamInterceptor %s %t %t", info.FullMethod, info.IsClientStream, info.IsServerStream)

	ctx := ss.Context()
	if err := authorize(ctx); err != nil {
		log.Infof("AUTH FAILED streamInterceptor %s %t %t", info.FullMethod, info.IsClientStream, info.IsServerStream)
		return err
	}

	log.Infof("AUTH OK streamInterceptor %s %t %t", info.FullMethod, info.IsClientStream, info.IsServerStream)

	// Continue execution of handler
	return handler(srv, ss)
}

// In case of failed authorization, the interceptor blocks execution of the handler and returns an error.
// type grpc.StreamClientInterceptor
func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Infof("unaryInterceptor %s", info.FullMethod)

	// Skip authorize when GetJWT is requested
	//if info.FullMethod != "/proto.EventStoreService/GetJWT" {
	//	if err := authorize(ctx); err != nil {
	//		return nil, err
	//	}
	//}

	if err := authorize(ctx); err != nil {
		return nil, err
	}

	// Continue execution of handler
	return handler(ctx, req)
}

// authorize ensures a valid token exists within a request's metadata and authorizes the token received from Metadata
func authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMissingMetadata
	}

	log.Infof("Dump Metadata ---")
	for key, strings := range md {
		log.Infof("[%s]=", key)
		for _, str := range strings {
			log.Infof("    %s", str)
		}
	}
	log.Infof("End Dump Metadata ---")

	authMetadata, ok := md["authorization"]
	if !ok {
		return errMissingToken
	}

	token := strings.TrimPrefix(authMetadata[0], "Bearer ")

	log.Infof("Bearer %s", token)
	err := validateToken(token)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	return nil
}

// valid validates the authorization.
func validateToken(_token string) error {

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(_token, func(token *jwt.Token) (interface{}, error) {

		// What also is used for JWT token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			// Do now support it
		}

		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			// Do now support it
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			// Support it
			// Return RSA Public Key

			rsaPK, err := jwt.ParseRSAPublicKeyFromPEM(jwtPublicKey)
			if err != nil {
				return nil, err
			}

			return rsaPK, nil
		}

		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	if err == nil {
		log.Infof("validateTokenOk")
	} else {
		log.Infof("validateTokenErr %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Infof("Looking for claims")
		for key, value := range claims {
			log.Infof("[%s]=%s", key, value)
		}
	}

	if token.Valid {
		return nil
	} else {
		return errInvalidToken
	}
}

func SetupOAuth(jwtPublicKeyFile string) ([]grpc.ServerOption, error) {
	opts := []grpc.ServerOption{
		// Add an interceptor for all unary RPCs.
		grpc.UnaryInterceptor(unaryInterceptor),

		// Add an interceptor for all stream RPCs.
		grpc.StreamInterceptor(streamInterceptor),
	}

	var err error
	jwtPublicKey, err = ioutil.ReadFile(jwtPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to access Public Key File '%s' file", jwtPublicKeyFile)
	}

	publicKeyPrefix := "-----BEGIN PUBLIC KEY-----"
	publicKeySuffix := "-----END PUBLIC KEY-----"
	if !strings.HasPrefix(string(jwtPublicKey), publicKeyPrefix) ||
		(!strings.HasSuffix(string(jwtPublicKey), publicKeySuffix) &&
			!strings.HasSuffix(string(jwtPublicKey), publicKeySuffix+"\n")) {
		return nil, fmt.Errorf("%s file must contain public key enclosed in %s %s",
			jwtPublicKeyFile, publicKeyPrefix, publicKeySuffix)
	}

	return opts, nil
}
