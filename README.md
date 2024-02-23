<strong>How to create a token?</strong><br />

    
	cre, err := authjwt.NewCredential(3600, secretkey, nil)

	if err != nil {
		return "", err
	}

	token, err := cre.CreateToken(username)

	if err != nil {
		return "", err
	}