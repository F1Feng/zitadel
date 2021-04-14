package management

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/eventstore/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func appFromModel(app *proj_model.Application) *management.Application {
	creationDate, err := ptypes.TimestampProto(app.CreationDate)
	logging.Log("GRPC-iejs3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(app.ChangeDate)
	logging.Log("GRPC-di7rw").OnError(err).Debug("unable to parse timestamp")

	return &management.Application{
		Id:           app.AppID,
		State:        appStateFromModel(app.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         app.Name,
		Sequence:     app.Sequence,
		AppConfig:    appConfigFromModel(app),
	}
}

func appConfigFromModel(app *proj_model.Application) management.AppConfig {
	if app.Type == proj_model.AppTypeOIDC {
		return &management.Application_OidcConfig{
			OidcConfig: oidcConfigFromModel(app.OIDCConfig),
		}
	}
	if app.Type == proj_model.AppTypeAPI {
		return &management.Application_ApiConfig{
			ApiConfig: apiConfigFromModel(app.APIConfig),
		}
	}
	return nil
}

func oidcConfigFromModel(config *proj_model.OIDCConfig) *management.OIDCConfig {
	return &management.OIDCConfig{
		RedirectUris:             config.RedirectUris,
		ResponseTypes:            oidcResponseTypesFromModel(config.ResponseTypes),
		GrantTypes:               oidcGrantTypesFromModel(config.GrantTypes),
		ApplicationType:          oidcApplicationTypeFromModel(config.ApplicationType),
		ClientId:                 config.ClientID,
		ClientSecret:             config.ClientSecretString,
		AuthMethodType:           oidcAuthMethodTypeFromModel(config.AuthMethodType),
		PostLogoutRedirectUris:   config.PostLogoutRedirectUris,
		Version:                  oidcVersionFromModel(config.OIDCVersion),
		NoneCompliant:            config.Compliance.NoneCompliant,
		ComplianceProblems:       complianceProblemsToLocalizedMessages(config.Compliance.Problems),
		DevMode:                  config.DevMode,
		AccessTokenType:          oidcTokenTypeFromModel(config.AccessTokenType),
		AccessTokenRoleAssertion: config.AccessTokenRoleAssertion,
		IdTokenRoleAssertion:     config.IDTokenRoleAssertion,
		IdTokenUserinfoAssertion: config.IDTokenUserinfoAssertion,
		ClockSkew:                durationpb.New(config.ClockSkew),
	}
}

func apiConfigFromModel(config *proj_model.APIConfig) *management.APIConfig {
	return &management.APIConfig{
		ClientId:       config.ClientID,
		ClientSecret:   config.ClientSecretString,
		AuthMethodType: apiAuthMethodTypeFromModel(config.AuthMethodType),
	}
}

func oidcConfigFromApplicationViewModel(app *proj_model.ApplicationView) *management.OIDCConfig {
	return &management.OIDCConfig{
		RedirectUris:             app.OIDCRedirectUris,
		ResponseTypes:            oidcResponseTypesFromModel(app.OIDCResponseTypes),
		GrantTypes:               oidcGrantTypesFromModel(app.OIDCGrantTypes),
		ApplicationType:          oidcApplicationTypeFromModel(app.OIDCApplicationType),
		ClientId:                 app.OIDCClientID,
		AuthMethodType:           oidcAuthMethodTypeFromModel(app.OIDCAuthMethodType),
		PostLogoutRedirectUris:   app.OIDCPostLogoutRedirectUris,
		Version:                  oidcVersionFromModel(app.OIDCVersion),
		NoneCompliant:            app.NoneCompliant,
		ComplianceProblems:       complianceProblemsToLocalizedMessages(app.ComplianceProblems),
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeFromModel(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IdTokenRoleAssertion:     app.IDTokenRoleAssertion,
		IdTokenUserinfoAssertion: app.IDTokenUserinfoAssertion,
		ClockSkew:                durationpb.New(app.ClockSkew),
	}
}

func apiConfigFromApplicationViewModel(app *proj_model.ApplicationView) *management.APIConfig {
	return &management.APIConfig{
		ClientId:       app.OIDCClientID,
		AuthMethodType: apiAuthMethodTypeFromModel(proj_model.APIAuthMethodType(app.OIDCAuthMethodType)),
	}
}

func complianceProblemsToLocalizedMessages(problems []string) []*message.LocalizedMessage {
	converted := make([]*message.LocalizedMessage, len(problems))
	for i, p := range problems {
		converted[i] = message.NewLocalizedMessage(p)
	}
	return converted

}

func oidcAppCreateToModel(app *management.OIDCApplicationCreate) *proj_model.Application {
	return &proj_model.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		Name: app.Name,
		Type: proj_model.AppTypeOIDC,
		OIDCConfig: &proj_model.OIDCConfig{
			OIDCVersion:              oidcVersionToModel(app.Version),
			RedirectUris:             app.RedirectUris,
			ResponseTypes:            oidcResponseTypesToModel(app.ResponseTypes),
			GrantTypes:               oidcGrantTypesToModel(app.GrantTypes),
			ApplicationType:          oidcApplicationTypeToModel(app.ApplicationType),
			AuthMethodType:           oidcAuthMethodTypeToModel(app.AuthMethodType),
			PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
			DevMode:                  app.DevMode,
			AccessTokenType:          oidcTokenTypeToModel(app.AccessTokenType),
			AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
			IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
			IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
			ClockSkew:                app.ClockSkew.AsDuration(),
		},
	}
}

func apiAppCreateToModel(app *management.APIApplicationCreate) *proj_model.Application {
	return &proj_model.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		Name: app.Name,
		Type: proj_model.AppTypeAPI,
		APIConfig: &proj_model.APIConfig{
			AuthMethodType: apiAuthMethodTypeToModel(app.AuthMethodType),
		},
	}
}

func appUpdateToModel(app *management.ApplicationUpdate) *proj_model.Application {
	return &proj_model.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID: app.Id,
		Name:  app.Name,
	}
}

func oidcConfigUpdateToModel(app *management.OIDCConfigUpdate) *proj_model.OIDCConfig {
	return &proj_model.OIDCConfig{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:                    app.ApplicationId,
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            oidcResponseTypesToModel(app.ResponseTypes),
		GrantTypes:               oidcGrantTypesToModel(app.GrantTypes),
		ApplicationType:          oidcApplicationTypeToModel(app.ApplicationType),
		AuthMethodType:           oidcAuthMethodTypeToModel(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeToModel(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
		ClockSkew:                app.ClockSkew.AsDuration(),
	}
}

func apiConfigUpdateToModel(app *management.APIConfigUpdate) *proj_model.APIConfig {
	return &proj_model.APIConfig{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:          app.ApplicationId,
		AuthMethodType: apiAuthMethodTypeToModel(app.AuthMethodType),
	}
}

func applicationSearchRequestsToModel(request *management.ApplicationSearchRequest) *proj_model.ApplicationSearchRequest {
	return &proj_model.ApplicationSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: applicationSearchQueriesToModel(request.ProjectId, request.Queries),
	}
}

func applicationSearchQueriesToModel(projectID string, queries []*management.ApplicationSearchQuery) []*proj_model.ApplicationSearchQuery {
	converted := make([]*proj_model.ApplicationSearchQuery, len(queries)+1)
	for i, q := range queries {
		converted[i] = applicationSearchQueryToModel(q)
	}
	converted[len(queries)] = &proj_model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Method: model.SearchMethodEquals, Value: projectID}

	return converted
}

func applicationSearchQueryToModel(query *management.ApplicationSearchQuery) *proj_model.ApplicationSearchQuery {
	return &proj_model.ApplicationSearchQuery{
		Key:    applicationSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func applicationSearchKeyToModel(key management.ApplicationSearchKey) proj_model.AppSearchKey {
	switch key {
	case management.ApplicationSearchKey_APPLICATIONSEARCHKEY_APP_NAME:
		return proj_model.AppSearchKeyName
	default:
		return proj_model.AppSearchKeyUnspecified
	}
}

func applicationSearchResponseFromModel(response *proj_model.ApplicationSearchResponse) *management.ApplicationSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Lp06f").OnError(err).Debug("unable to parse timestamp")
	return &management.ApplicationSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            applicationViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func applicationViewsFromModel(apps []*proj_model.ApplicationView) []*management.ApplicationView {
	converted := make([]*management.ApplicationView, len(apps))
	for i, app := range apps {
		converted[i] = applicationViewFromModel(app)
	}
	return converted
}

func applicationViewFromModel(application *proj_model.ApplicationView) *management.ApplicationView {
	creationDate, err := ptypes.TimestampProto(application.CreationDate)
	logging.Log("GRPC-lo9sw").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(application.ChangeDate)
	logging.Log("GRPC-8uwsd").OnError(err).Debug("unable to parse timestamp")

	converted := &management.ApplicationView{
		Id:           application.ID,
		State:        appStateFromModel(application.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         application.Name,
		Sequence:     application.Sequence,
	}
	if application.IsOIDC {
		converted.AppConfig = &management.ApplicationView_OidcConfig{
			OidcConfig: oidcConfigFromApplicationViewModel(application),
		}
	} else {
		converted.AppConfig = &management.ApplicationView_ApiConfig{
			ApiConfig: apiConfigFromApplicationViewModel(application),
		}
	}
	return converted
}

func appStateFromModel(state proj_model.AppState) management.AppState {
	switch state {
	case proj_model.AppStateActive:
		return management.AppState_APPSTATE_ACTIVE
	case proj_model.AppStateInactive:
		return management.AppState_APPSTATE_INACTIVE
	default:
		return management.AppState_APPSTATE_UNSPECIFIED
	}
}

func oidcResponseTypesToModel(responseTypes []management.OIDCResponseType) []proj_model.OIDCResponseType {
	if responseTypes == nil || len(responseTypes) == 0 {
		return []proj_model.OIDCResponseType{proj_model.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]proj_model.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case management.OIDCResponseType_OIDCRESPONSETYPE_CODE:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeCode
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeIDToken
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeIDTokenToken
		}
	}

	return oidcResponseTypes
}

func oidcResponseTypesFromModel(responseTypes []proj_model.OIDCResponseType) []management.OIDCResponseType {
	oidcResponseTypes := make([]management.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case proj_model.OIDCResponseTypeCode:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_CODE
		case proj_model.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN
		case proj_model.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN
		}
	}

	return oidcResponseTypes
}

func oidcGrantTypesToModel(grantTypes []management.OIDCGrantType) []proj_model.OIDCGrantType {
	if grantTypes == nil || len(grantTypes) == 0 {
		return []proj_model.OIDCGrantType{proj_model.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]proj_model.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeAuthorizationCode
		case management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeImplicit
		case management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeRefreshToken
		}
	}
	return oidcGrantTypes
}

func oidcGrantTypesFromModel(grantTypes []proj_model.OIDCGrantType) []management.OIDCGrantType {
	oidcGrantTypes := make([]management.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case proj_model.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE
		case proj_model.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT
		case proj_model.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToModel(appType management.OIDCApplicationType) proj_model.OIDCApplicationType {
	switch appType {
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB:
		return proj_model.OIDCApplicationTypeWeb
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT:
		return proj_model.OIDCApplicationTypeUserAgent
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE:
		return proj_model.OIDCApplicationTypeNative
	}
	return proj_model.OIDCApplicationTypeWeb
}

func oidcVersionToModel(version management.OIDCVersion) proj_model.OIDCVersion {
	switch version {
	case management.OIDCVersion_OIDCV1_0:
		return proj_model.OIDCVersionV1
	}
	return proj_model.OIDCVersionV1
}

func oidcApplicationTypeFromModel(appType proj_model.OIDCApplicationType) management.OIDCApplicationType {
	switch appType {
	case proj_model.OIDCApplicationTypeWeb:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	case proj_model.OIDCApplicationTypeUserAgent:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT
	case proj_model.OIDCApplicationTypeNative:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE
	default:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	}
}

func oidcAuthMethodTypeToModel(authType management.OIDCAuthMethodType) proj_model.OIDCAuthMethodType {
	switch authType {
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC:
		return proj_model.OIDCAuthMethodTypeBasic
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST:
		return proj_model.OIDCAuthMethodTypePost
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE:
		return proj_model.OIDCAuthMethodTypeNone
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_PRIVATE_KEY_JWT:
		return proj_model.OIDCAuthMethodTypePrivateKeyJWT
	default:
		return proj_model.OIDCAuthMethodTypeBasic
	}
}

func apiAuthMethodTypeToModel(authType management.APIAuthMethodType) proj_model.APIAuthMethodType {
	switch authType {
	case management.APIAuthMethodType_APIAUTHMETHODTYPE_BASIC:
		return proj_model.APIAuthMethodTypeBasic
	case management.APIAuthMethodType_APIAUTHMETHODTYPE_PRIVATE_KEY_JWT:
		return proj_model.APIAuthMethodTypePrivateKeyJWT
	default:
		return proj_model.APIAuthMethodTypeBasic
	}
}

func oidcAuthMethodTypeFromModel(authType proj_model.OIDCAuthMethodType) management.OIDCAuthMethodType {
	switch authType {
	case proj_model.OIDCAuthMethodTypeBasic:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	case proj_model.OIDCAuthMethodTypePost:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST
	case proj_model.OIDCAuthMethodTypeNone:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE
	case proj_model.OIDCAuthMethodTypePrivateKeyJWT:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_PRIVATE_KEY_JWT
	default:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	}
}

func apiAuthMethodTypeFromModel(authType proj_model.APIAuthMethodType) management.APIAuthMethodType {
	switch authType {
	case proj_model.APIAuthMethodTypeBasic:
		return management.APIAuthMethodType_APIAUTHMETHODTYPE_BASIC
	case proj_model.APIAuthMethodTypePrivateKeyJWT:
		return management.APIAuthMethodType_APIAUTHMETHODTYPE_PRIVATE_KEY_JWT
	default:
		return management.APIAuthMethodType_APIAUTHMETHODTYPE_BASIC
	}
}

func oidcTokenTypeToModel(tokenType management.OIDCTokenType) proj_model.OIDCTokenType {
	switch tokenType {
	case management.OIDCTokenType_OIDCTokenType_Bearer:
		return proj_model.OIDCTokenTypeBearer
	case management.OIDCTokenType_OIDCTokenType_JWT:
		return proj_model.OIDCTokenTypeJWT
	default:
		return proj_model.OIDCTokenTypeBearer
	}
}

func oidcTokenTypeFromModel(tokenType proj_model.OIDCTokenType) management.OIDCTokenType {
	switch tokenType {
	case proj_model.OIDCTokenTypeBearer:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	case proj_model.OIDCTokenTypeJWT:
		return management.OIDCTokenType_OIDCTokenType_JWT
	default:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	}
}

func oidcVersionFromModel(version proj_model.OIDCVersion) management.OIDCVersion {
	switch version {
	case proj_model.OIDCVersionV1:
		return management.OIDCVersion_OIDCV1_0
	default:
		return management.OIDCVersion_OIDCV1_0
	}
}

func appChangesToResponse(response *proj_model.ApplicationChanges, offset uint64, limit uint64) (_ *management.Changes) {
	return &management.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: appChangesToMgtAPI(response),
	}
}

func appChangesToMgtAPI(changes *proj_model.ApplicationChanges) (_ []*management.Change) {
	result := make([]*management.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &management.Change{
			ChangeDate: change.ChangeDate,
			EventType:  message.NewLocalizedEventType(change.EventType),
			Sequence:   change.Sequence,
			Editor:     change.ModifierName,
			EditorId:   change.ModifierId,
			Data:       data,
		}
	}

	return result
}

func clientKeyViewsFromModel(keys ...*key_model.AuthNKeyView) []*management.ClientKeyView {
	keyViews := make([]*management.ClientKeyView, len(keys))
	for i, key := range keys {
		keyViews[i] = clientKeyViewFromModel(key)
	}
	return keyViews
}

func clientKeyViewFromModel(key *key_model.AuthNKeyView) *management.ClientKeyView {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-DAs2t").OnError(err).Debug("unable to parse timestamp")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-BDgh4").OnError(err).Debug("unable to parse timestamp")

	return &management.ClientKeyView{
		Id:             key.ID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		Type:           authNKeyTypeFromModel(key.Type),
	}
}

func addClientKeyToModel(key *management.AddClientKeyRequest) *proj_model.ClientKey {
	expirationDate := time.Time{}
	if key.ExpirationDate != nil {
		var err error
		expirationDate, err = ptypes.Timestamp(key.ExpirationDate)
		logging.Log("MANAG-Dgt42").OnError(err).Debug("unable to parse expiration date")
	}

	return &proj_model.ClientKey{
		ExpirationDate: expirationDate,
		Type:           authNKeyTypeToModel(key.Type),
		ApplicationID:  key.ApplicationId,
		ObjectRoot:     models.ObjectRoot{AggregateID: key.ProjectId},
	}
}

func addClientKeyFromModel(key *proj_model.ClientKey) *management.AddClientKeyResponse {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-FBzz4").OnError(err).Debug("unable to parse cretaion date")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-sag21").OnError(err).Debug("unable to parse cretaion date")

	detail, err := json.Marshal(struct {
		Type     string `json:"type"`
		KeyID    string `json:"keyId"`
		Key      string `json:"key"`
		AppID    string `json:"appId"`
		ClientID string `json:"clientId"`
	}{
		Type:     "application",
		KeyID:    key.KeyID,
		Key:      string(key.PrivateKey),
		AppID:    key.ApplicationID,
		ClientID: key.ClientID,
	})
	logging.Log("MANAG-adt42").OnError(err).Warn("unable to marshall key")

	return &management.AddClientKeyResponse{
		Id:             key.KeyID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		KeyDetails:     detail,
		Type:           authNKeyTypeFromModel(key.Type),
	}
}

func authNKeyTypeToModel(typ management.AuthNKeyType) key_model.AuthNKeyType {
	switch typ {
	case management.AuthNKeyType_AUTHNKEY_JSON:
		return key_model.AuthNKeyTypeJSON
	default:
		return key_model.AuthNKeyTypeNONE
	}
}

func authNKeyTypeFromModel(typ key_model.AuthNKeyType) management.AuthNKeyType {
	switch typ {
	case key_model.AuthNKeyTypeJSON:
		return management.AuthNKeyType_AUTHNKEY_JSON
	default:
		return management.AuthNKeyType_AUTHNKEY_UNSPECIFIED
	}
}

func clientKeySearchRequestToModel(req *management.ClientKeySearchRequest) *key_model.AuthNKeySearchRequest {
	return &key_model.AuthNKeySearchRequest{
		Offset: req.Offset,
		Limit:  req.Limit,
		Asc:    req.Asc,
		Queries: []*key_model.AuthNKeySearchQuery{
			{
				Key:    key_model.AuthNKeyObjectType,
				Method: model.SearchMethodEquals,
				Value:  key_model.AuthNKeyObjectTypeApplication,
			}, {
				Key:    key_model.AuthNKeyObjectID,
				Method: model.SearchMethodEquals,
				Value:  req.ApplicationId,
			},
		},
	}
}

func clientKeySearchResponseFromModel(req *key_model.AuthNKeySearchResponse) *management.ClientKeySearchResponse {
	viewTimestamp, err := ptypes.TimestampProto(req.Timestamp)
	logging.Log("MANAG-Sk9ds").OnError(err).Debug("unable to parse cretaion date")

	return &management.ClientKeySearchResponse{
		Offset:            req.Offset,
		Limit:             req.Limit,
		TotalResult:       req.TotalResult,
		ProcessedSequence: req.Sequence,
		ViewTimestamp:     viewTimestamp,
		Result:            clientKeyViewsFromModel(req.Result...),
	}
}