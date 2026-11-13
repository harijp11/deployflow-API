package clients

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "deployflow/auth-service/internal/services"
)

type UserServiceClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
    return &UserServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

func (client *UserServiceClient) CreateUser(
    ctx context.Context,
    profile services.UserProfile,
) (uint, error) {
    body, err := json.Marshal(profile)
    if err != nil {
        return 0, err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        client.baseURL+"/api/v1/users",
        bytes.NewBuffer(body),
    )
    if err != nil {
        return 0, err
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := client.httpClient.Do(req)
    if err != nil {
        return 0, fmt.Errorf("call user service: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return 0, fmt.Errorf("user service returned status %d", resp.StatusCode)
    }

    var result struct {
        Data struct {
            ID uint `json:"id"`
        } `json:"data"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, fmt.Errorf("decode user service response: %w", err)
    }

    if result.Data.ID == 0 {
        return 0, fmt.Errorf("user service returned an invalid user ID")
    }

    return result.Data.ID, nil
}

var _ services.UserCreator = (*UserServiceClient)(nil)