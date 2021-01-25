package usecases

import (
	"context"
	"errors"

	"assignment/constants"
	"assignment/entity"
)

func CreateGithubToken(reader Reader, githubToken GithubToken) *AuthAccessToken {
	return &AuthAccessToken{
		reader: reader,
		gt:     githubToken,
	}
}

type GithubToken interface {
	AccessToken(context.Context, entity.OAuthRedirectValues) error
}

type AuthAccessToken struct {
	reader Reader
	gt     GithubToken
}

func (aat *AuthAccessToken) Authorize(ctx context.Context) error {
	oAuthValues, err := aat.reader.Get(ctx)
	if err != nil {
		return err
	}
	code := ctx.Value(constants.CodeQueryParam).(string)
	if code == "" {
		return errors.New("Code not found in context")
	}
	return aat.gt.AccessToken(ctx, oAuthValues)

}
