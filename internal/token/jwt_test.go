package token

import "testing"

func TestJWTManager_GenerateAndParse(t *testing.T) {
	secret := "super_secret_key_for_testing_61123"
	jwtManager := NewJWTManager(secret)

	userID := int64(1)

	token, err := jwtManager.Generate(userID)
	if err != nil {
		t.Fatalf("token generate error: %v", err)
	}

	parsedTokenID, err := jwtManager.Parse(token)
	if err != nil {
		t.Fatalf("token parse error: %v", err)
	}

	if parsedTokenID != userID {
		t.Fatalf("Parse(%v) expected id: %d; get: %d", token, userID, parsedTokenID)
	}
}

func TestJWTManager_Parse_InvalidToken(t *testing.T) {
	manager1 := NewJWTManager("secret1")
	manager2 := NewJWTManager("secret2")

	tokenString, _ := manager1.Generate(1)

	_, err := manager2.Parse(tokenString)
	if err == nil {
		t.Errorf("expected error with another jwt manager, but haven`t error")
	}
}
