<strong>Version: v1.19.0</strong><br />
<strong>Adding Bearer authentication with Gin</strong><br/>
 ```go
 import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rafaelsouzaribeiro/jwt-auth/pkg/middleware"
)

router := gin.Default()
cre, err := middleware.NewCredential(3600, "secretkey", nil)

if err != nil {
	panic(err)
}

router.POST("/publish", cre.AuthMiddlewareGin(), func(c *gin.Context) {
	var json struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": json.Message,
	})
})

router.Run(":8080")
 ```

<strong>Version: v1.15.0 or latest</strong><br />
<strong>Updating folder: pkg/middleware</strong><br /> 
 ```go
import jwtauth "github.com/rafaelsouzaribeiro/jwt-auth/pkg/middleware"
 ```

<strong>Version: v1.14.0 or latest</strong><br />
<strong>ExtractClaims method</strong><br /> 
 ```go

	cre, err := jwtauth.NewCredential(3600, "secretkey", nil)

	if err != nil {
		panic(err)
	}	

	claims, err = cre.ExtractClaims(token)

	if err != nil {
		panic(err)
	}

	println(fmt.Sprint(claims["lastname"]), fmt.Sprint(claims["firstname"]))
```
<br /> 
<strong>CreateToken</strong><br />
	
```go 	
	// pass multiple data
	claims := map[string]interface{}{
		"username": username,
		 // ... other claims
	}

	token, err := cre.CreateToken(claims)
```

<strong>Methods for http getting the bearer token and validating</strong><br />      

  ```go
 package main

import (
	"fmt"
	"net/http"

	jwtauth "github.com/rafaelsouzaribeiro/jwt-auth"
)

func main() {

	mx := http.NewServeMux()
	cre, err := jwtauth.NewCredential(3600, "secretkey", nil)

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

<br/>
<strong>How to add an interceptor in grpc?</strong><br />  
<strong>This example takes Bearer Token Authentication and skips token validation for functions login,loginAdm</strong><br />  

  ```go 
	c, errs := jwtauth.NewCredential(3600, secretkey, []string{"login", "loginAdm"})

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
	cre, err := jwtauth.NewCredential(3600, secretkey, nil)

	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"username": username,
		 // ... other claims
	}

	token, err := cre.CreateToken(claims)

	if err != nil {
		return "", err
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(c.JwtUnaryInterceptor(token)),
		grpc.StreamInterceptor(c.JwtStreamInterceptor(token)),
	)

  ```
 
