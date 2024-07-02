package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emsi-zero/auth_ir/internal/conf"
	"github.com/emsi-zero/auth_ir/internal/models"
	"github.com/emsi-zero/auth_ir/internal/utilities"
	"github.com/pkg/errors"
)

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func isAdmin(u *models.User, config *conf.GlobalConfiguration) bool {
	return config.JWT.Aud == u.Aud && u.HasRole(config.JWT.AdminGroupName)
}

func (a *API) requestAud(ctx context.Context, r *http.Request) string {
	config := a.config
	// First check for an audience in the header
	if aud := r.Header.Get(audHeaderName); aud != "" {
		return aud
	}

	// Then check the token
	claims := getClaims(ctx)
	if claims != nil && claims.Audience != "" {
		return claims.Audience
	}

	// Finally, return the default if none of the above methods are successful
	return config.JWT.Aud
}

func isStringInSlice(checkValue string, list []string) bool {
	for _, val := range list {
		if val == checkValue {
			return true
		}
	}
	return false
}

// getBodyBytes returns a byte array of the request's Body.
func getBodyBytes(req *http.Request) ([]byte, error) {
	return utilities.GetBodyBytes(req)
}

type RequestParams interface {
	AdminUserParams |
		CreateSSOProviderParams |
		EnrollFactorParams |
		GenerateLinkParams |
		IdTokenGrantParams |
		InviteParams |
		OtpParams |
		PKCEGrantParams |
		PasswordGrantParams |
		RecoverParams |
		RefreshTokenGrantParams |
		ResendConfirmationParams |
		SignupParams |
		SingleSignOnParams |
		SmsParams |
		UserUpdateParams |
		VerifyFactorParams |
		VerifyParams |
		adminUserUpdateFactorParams |
		struct {
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
}

// retrieveRequestParams is a generic method that unmarshals the request body into the params struct provided
func retrieveRequestParams[A RequestParams](r *http.Request, params *A) error {
	body, err := getBodyBytes(r)
	if err != nil {
		return internalServerError("Could not read body into byte slice").WithInternalError(err)
	}
	if err := json.Unmarshal(body, params); err != nil {
		return badRequestError(ErrorCodeBadJSON, "Could not parse request body as JSON: %v", err)
	}
	return nil
}
