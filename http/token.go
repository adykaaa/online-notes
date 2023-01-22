package http

symmetricKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // Must be 32 bytes
now := time.Now()
exp := now.Add(24 * time.Hour)
nbt := now

jsonToken := paseto.JSONToken{
        Audience:   "test",
        Issuer:     "test_service",
        Jti:        "123",
        Subject:    "test_subject",
        IssuedAt:   now,
        Expiration: exp,
        NotBefore:  nbt,
        }
// Add custom claim    to the token    
jsonToken.Set("data", "this is a signed message")
footer := "some footer"

// Encrypt data
token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
// token = "v2.local.E42A2iMY9SaZVzt-WkCi45_aebky4vbSUJsfG45OcanamwXwieieMjSjUkgsyZzlbYt82miN1xD-X0zEIhLK_RhWUPLZc9nC0shmkkkHS5Exj2zTpdNWhrC5KJRyUrI0cupc5qrctuREFLAvdCgwZBjh1QSgBX74V631fzl1IErGBgnt2LV1aij5W3hw9cXv4gtm_jSwsfee9HZcCE0sgUgAvklJCDO__8v_fTY7i_Regp5ZPa7h0X0m3yf0n4OXY9PRplunUpD9uEsXJ_MTF5gSFR3qE29eCHbJtRt0FFl81x-GCsQ9H9701TzEjGehCC6Bhw.c29tZSBmb290ZXI"

// Decrypt data
var newJsonToken paseto.JSONToken
var newFooter string
err := paseto.Decrypt(token, symmetricKey, &newJsonToken, &newFooter)