<strong>How to create a token?</strong><br />

    
	cre, err := authjwt.NewCredential(3600, secretkey, nil)

	if err != nil {
		panic(err)
	}

	token, err := cre.CreateToken(username)

	if err != nil {
		panic(err)
	}

<strong>How to add an interceptor in grpc?</strong><br />  
<strong>This example takes Bearer Token Authentication and skips token validation for functions login,loginAdm</strong><br />  

	c, errs := authjwt.NewCredential(3600, secretkey, []string{"login", "loginAdm"})

	if err != nil {
		panic(errs)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(c.UnaryInterceptorBearer),
		grpc.StreamInterceptor(c.StreamInterceptorBearer),
	)

<strong>How to add an interceptor in grpc? passing the token as a parameter</strong><br />   

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