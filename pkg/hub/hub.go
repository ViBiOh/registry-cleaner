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
	req request.Request
}

type tagResponse struct {
	Count uint64 `json:"count"`
}

type tagResult struct {
	Name string `json:"name"`
}

// New creates new App from Config
func New(ctx context.Context, username, password string) (App, error) {
	jwt, err := login(ctx, username, password)
	if err != nil {
		return App{}, err
	}

	return App{
		req: request.New().URL("https://hub.docker.com/v2/repositories").Header("Authorization", fmt.Sprintf("JWT %s", jwt)),
	}, nil
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

	resp, err := a.req.Method(http.MethodGet).Path(fmt.Sprintf("/%s/tags/?page_size=1", image)).Send(ctx, nil)
	if err != nil {
		return err
	}

	var tagsCount tagResponse
	if err = httpjson.Read(resp, &tagsCount); err != nil {
		return fmt.Errorf("unable to parse tags' count: %s", err)
	}

	for page := 1; page <= int(math.Ceil(float64(tagsCount.Count)/pageSize)); page++ {
		resp, err := a.req.Method(http.MethodGet).Path(fmt.Sprintf("/%s/tags/?page_size=%d&page=%d", image, pageSize, page)).Send(ctx, nil)
		if err != nil {
			return err
		}

		if err = httpjson.Stream(resp.Body, tags, "results", false); err != nil {
			return fmt.Errorf("unable to parse json: %s", err)
		}
	}

	close(tags)
	<-done

	return nil
}

// Delete a tag
func (a App) Delete(ctx context.Context, image, tag string) error {
	resp, err := a.req.Method(http.MethodDelete).Path(fmt.Sprintf("/%s/tags/%s/", image, tag)).Send(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to delete tag: %s", err)
	}

	if err = request.DiscardBody(resp.Body); err != nil {
		return fmt.Errorf("unable to discard body: %s", err)
	}

	return nil
}

func login(ctx context.Context, username, password string) (string, error) {
	loginPayload := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := request.Post("https://hub.docker.com/v2/users/login/").AcceptJSON().JSON(context.Background(), loginPayload)
	if err != nil {
		return "", fmt.Errorf("unable to login to Docker Hub: %s", err)
	}

	var output map[string]string
	if err = httpjson.Read(resp, &output); err != nil {
		return "", fmt.Errorf("unable to parse login response from Docker Hub: %s", err)
	}

	return output["token"], nil
}
