

<strong>Version: v1.14.0</strong><br />
<strong>update repository name</strong><br /> 
<p>
update repository name from auth-jwt to jwt-auth
</p>

<strong>Version: v1.13.0</strong><br />
<strong>Added ExtractClaims method</strong><br /> 
 ```go

	cre, err := authjwt.NewCredential(1, "secretkey", nil)

	if err != nil {
		panic(err)
	}

	claims, err = cre.ExtractClaims(token)

	if err != nil {
		panic(err)
	}

	println(fmt.Sprint(claims["lastname"]), fmt.Sprint(claims["firstname"]))
```
<strong>Version: v1.12.0 or latest</strong><br />
<strong>Methods for http getting the bearer token and validating</strong><br />      

  ```go
 package main

import (
	"fmt"
	"net/http"

	authjwt "github.com/rafaelsouzaribeiro/auth-jwt"
)

func main() {

	mx := http.NewServeMux()
	cre, err := authjwt.NewCredential(3600, "secretkey", nil)

	if err != nil {
		panic(err)
	}

	// Protected routes
	mx.HandleFunc("/route1", cre.AuthMiddleware(rota1Handler))
	mx.HandleFunc("/route2", cre.AuthMiddleware(rota2Handler))

	// Public route
	mx.HandleFunc("/public-route", rotaPublicaHandler)

	http.ListenAndServe(":8080", mx)
}

func rota1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Token-protected Route 1")
}

func rota2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Token-protected Route 2")
}

func rotaPublicaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Public route accessible without token")
}


```

<strong>Version: v1.11.0 or latest</strong><br />
<strong>Update of the new CreateToken</strong><br />
	
```go 	
	// pass multiple data
	claims := map[string]interface{}{
		"username": username,
		 // ... other claims
	}

	token, err := cre.CreateToken(claims)
```

<strong>Version: v1.10.0</strong><br />
<strong>How to create a token?</strong><br />

  ```go 
	cre, err := authjwt.NewCredential(3600, secretkey, nil)

	if err != nil {
		panic(err)
	}

	token, err := cre.CreateToken(username)

	if err != nil {
		panic(err)
	}
```

<strong>How to add an interceptor in grpc?</strong><br />  
<strong>This example takes Bearer Token Authentication and skips token validation for functions login,loginAdm</strong><br />  

  ```go 
	c, errs := authjwt.NewCredential(3600, secretkey, []string{"login", "loginAdm"})

	if err != nil {
		panic(errs)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(c.UnaryInterceptorBearer),
		grpc.StreamInterceptor(c.StreamInterceptorBearer),
	)
```
<strong>How to add an interceptor in grpc? passing the token as a parameter</strong><br />   

  ```go
	cre, err := authjwt.NewCredential(3600, secretkey, nil)

	if err != nil {
		return "", err
	}

	token, err := cre.CreateToken(username)

	if err != nil {
		return "", err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(c.JwtUnaryInterceptor(token)),
		grpc.StreamInterceptor(c.JwtStreamInterceptor(token)),
	)

  ```
 
<strong>Methods for http</strong><br />      

  ```go
    token, err := cre.CreateToken(username)
    err := cre.VerifyToken(token)
    	if cre.TokenExpired(token) {
			
		}
		  ```