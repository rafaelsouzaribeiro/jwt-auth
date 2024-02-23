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