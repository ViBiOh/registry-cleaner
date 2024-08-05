package hub

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"runtime"

	"github.com/ViBiOh/httputils/v4/pkg/httpjson"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

const pageSize = 100

// App of package
type App struct {
	owner string
	req   request.Request
}

type listResult struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type tagResponse struct {
	Count uint64 `json:"count"`
}

type tagResult struct {
	Name string `json:"name"`
}

// New creates new App from Config
func New(ctx context.Context, username, password, owner string) (App, error) {
	jwt, err := login(ctx, username, password)
	if err != nil {
		return App{}, err
	}

	if len(owner) == 0 {
		owner = username
	}

	return App{
		req:   request.New().URL("https://hub.docker.com/v2").Header("Authorization", fmt.Sprintf("JWT %s", jwt)),
		owner: owner,
	}, nil
}

// Repositories list repositories
func (a App) Repositories(ctx context.Context) ([]string, error) {
	resp, err := a.req.Method(http.MethodGet).Path("/users/%s/repositories/?page_size=100", a.owner).Send(ctx, nil)
	if err != nil {
		return nil, err
	}

	repositories, err := httpjson.Read[[]listResult](resp)
	if err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}

	output := make([]string, len(repositories))
	for i, repository := range repositories {
		output[i] = fmt.Sprintf("%s/%s", repository.Namespace, repository.Name)
	}

	return output, nil
}

// Tags list tags for a given image
func (a App) Tags(ctx context.Context, image string, handler func(string)) error {
	done := make(chan struct{})
	tags := make(chan tagResult, runtime.NumCPU())

	go func() {
		defer close(done)

		for tag := range tags {
			handler(tag.Name)
		}
	}()

	resp, err := a.req.Method(http.MethodGet).Path("/repositories/%s/tags/?page_size=1", image).Send(ctx, nil)
	if err != nil {
		return err
	}

	tagsCount, err := httpjson.Read[tagResponse](resp)
	if err != nil {
		return fmt.Errorf("parse tags' count: %w", err)
	}

	for page := 1; page <= int(math.Ceil(float64(tagsCount.Count)/pageSize)); page++ {
		resp, err := a.req.Method(http.MethodGet).Path("/repositories/%s/tags/?page_size=%d&page=%d", image, pageSize, page).Send(ctx, nil)
		if err != nil {
			return err
		}

		if err = httpjson.Stream(resp.Body, tags, "results", false); err != nil {
			return fmt.Errorf("parse json: %w", err)
		}

		if err := request.DiscardBody(resp.Body); err != nil {
			return fmt.Errorf("discard body: %w", err)
		}
	}

	close(tags)
	<-done

	return nil
}

// Delete a tag
func (a App) Delete(ctx context.Context, image, tag string) error {
	resp, err := a.req.Method(http.MethodDelete).Path("/repositories/%s/tags/%s/", image, tag).Send(ctx, nil)
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	if err = request.DiscardBody(resp.Body); err != nil {
		return fmt.Errorf("discard body: %w", err)
	}

	return nil
}

func login(ctx context.Context, username, password string) (string, error) {
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := request.Post("https://hub.docker.com/v2/users/login/").AcceptJSON().JSON(ctx, loginPayload)
	if err != nil {
		return "", fmt.Errorf("login to Docker Hub: %w", err)
	}

	output, err := httpjson.Read[map[string]string](resp)
	if err != nil {
		return "", fmt.Errorf("parse login response from Docker Hub: %w", err)
	}

	return output["token"], nil
}
