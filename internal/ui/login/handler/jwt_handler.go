package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caos/oidc/pkg/client/rp"
	"github.com/caos/oidc/pkg/oidc"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type jwtRequest struct {
	AuthRequestID string `schema:"authRequestID"`
	UserAgentID   string `schema:"userAgentID"`
}

func (l *Login) handleJWTRequest(w http.ResponseWriter, r *http.Request) {
	data := new(jwtRequest)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	id, err := base64.RawURLEncoding.DecodeString(data.UserAgentID)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, err := l.IDPConfigAesCrypto.DecryptString(id, l.IDPConfigAesCrypto.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.AuthRequestID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if idpConfig.IsOIDC {
		if err != nil {
			l.renderError(w, r, nil, err)
			return
		}
	}
	l.handleJWTExtraction(w, r, authReq, idpConfig)
}

func (l *Login) handleJWTExtraction(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, idpConfig *iam_model.IDPConfigView) {
	token, err := getToken(r)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	tokenClaims, err := validateToken(r.Context(), token, idpConfig)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	externalUser := l.mapTokenToLoginUser(&oidc.Tokens{IDTokenClaims: tokenClaims}, idpConfig)
	err = l.authRepo.CheckExternalUserLogin(r.Context(), authReq.ID, authReq.AgentID, externalUser, domain.BrowserInfoFromRequest(r))
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		}
		if !idpConfig.AutoRegister {
			l.renderExternalNotFoundOption(w, r, authReq, err)
			return
		}
		authReq, err = l.authRepo.AuthRequestByID(r.Context(), authReq.ID, authReq.AgentID)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, err)
			return
		}
		resourceOwner := l.getOrgID(authReq)
		orgIamPolicy, err := l.getOrgIamPolicy(r, resourceOwner)
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, err)
			return
		}
		user, externalIDP := l.mapExternalUserToLoginUser(orgIamPolicy, authReq.LinkingUsers[len(authReq.LinkingUsers)-1], idpConfig)
		err = l.authRepo.AutoRegisterExternalUser(setContext(r.Context(), resourceOwner), user, externalIDP, nil, authReq.ID, authReq.AgentID, resourceOwner, domain.BrowserInfoFromRequest(r))
		if err != nil {
			l.renderExternalNotFoundOption(w, r, authReq, err)
			return
		}
	}
	redirect, err := l.redirectToJWTCallback(authReq)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func (l *Login) redirectToJWTCallback(authReq *domain.AuthRequest) (string, error) {
	redirect, err := url.Parse(l.baseURL + EndpointJWTCallback)
	if err != nil {
		return "", err
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, authReq.ID)
	nonce, err := l.IDPConfigAesCrypto.Encrypt([]byte(authReq.AgentID))
	if err != nil {
		return "", err
	}
	q.Set(queryUserAgentID, base64.RawURLEncoding.EncodeToString(nonce))
	redirect.RawQuery = q.Encode()
	return redirect.String(), nil
}

func (l *Login) handleJWTCallback(w http.ResponseWriter, r *http.Request) {
	data := new(jwtRequest)
	err := l.getParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	id, err := base64.RawURLEncoding.DecodeString(data.UserAgentID)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	userAgentID, err := l.IDPConfigAesCrypto.DecryptString(id, l.IDPConfigAesCrypto.EncryptionKeyID())
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}
	authReq, err := l.authRepo.AuthRequestByID(r.Context(), data.AuthRequestID, userAgentID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	idpConfig, err := l.authRepo.GetIDPConfigByID(r.Context(), authReq.SelectedIDPConfigID)
	if err != nil {
		l.renderError(w, r, authReq, err)
		return
	}
	if idpConfig.IsOIDC {
		l.renderLogin(w, r, authReq, err)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func validateToken(ctx context.Context, token string, config *iam_model.IDPConfigView) (oidc.IDTokenClaims, error) {
	offset := 3 * time.Second
	maxAge := time.Hour
	claims := oidc.EmptyIDTokenClaims()
	payload, err := oidc.ParseToken(token, claims)
	if err != nil {
		return nil, err
	}

	if err := oidc.CheckSubject(claims); err != nil {
		return nil, err
	}

	if err = oidc.CheckIssuer(claims, config.JWTIssuer); err != nil {
		return nil, err
	}

	keySet := rp.NewRemoteKeySet(http.DefaultClient, config.JWTKeysEndpoint)
	if err = oidc.CheckSignature(ctx, token, payload, claims, nil, keySet); err != nil {
		return nil, err
	}

	if !claims.GetExpiration().IsZero() {
		if err = oidc.CheckExpiration(claims, offset); err != nil {
			return nil, err
		}
	}

	if err = oidc.CheckIssuedAt(claims, maxAge, offset); err != nil {
		return nil, err
	}
	return claims, nil
}

func getToken(r *http.Request) (string, error) {
	auth := r.Header.Get(http_util.Authorization)
	if auth == "" {
		auth = r.Header.Get("x-authorization")
	}
	if auth == "" {
		return "", errors.ThrowInvalidArgument(nil, "LOGIN-adh42", "Errors.AuthRequest.TokenNotFound")
	}
	return strings.TrimPrefix(auth, oidc.PrefixBearer), nil
}
