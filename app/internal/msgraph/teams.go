package msgraph

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// TeamsService reads Microsoft Teams channels and files via Microsoft Graph.
type TeamsService struct {
	client *Client
}

// NewTeamsService creates a Teams service backed by the given Graph client.
func NewTeamsService(client *Client) *TeamsService {
	return &TeamsService{client: client}
}

// Team represents a Microsoft Teams team.
type Team struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	WebURL      string `json:"webUrl"`
}

// Channel represents a Teams channel.
type Channel struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	WebURL      string `json:"webUrl"`
}

// TeamFile represents a file in a Teams channel's shared folder.
type TeamFile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WebURL         string `json:"webUrl"`
	Size           int64  `json:"size"`
	LastModified   string `json:"lastModifiedDateTime"`
	LastModifiedBy string
	IsFolder       bool
}

// ListTeams returns all teams the app has access to.
func (s *TeamsService) ListTeams() ([]Team, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	body, err := s.client.Get("/groups", map[string]string{
		"$filter": "resourceProvisioningOptions/Any(x:x eq 'Team')",
		"$select": "id,displayName,description",
		"$top":    "50",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value []Team `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	slog.Info("Teams listed", "count", len(resp.Value))
	return resp.Value, nil
}

// ListChannels returns channels for a specific team.
func (s *TeamsService) ListChannels(teamID string) ([]Channel, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	path := fmt.Sprintf("/teams/%s/channels", teamID)
	body, err := s.client.Get(path, map[string]string{
		"$select": "id,displayName,description,webUrl",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value []Channel `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// ListChannelFiles returns files in a Teams channel's shared folder.
func (s *TeamsService) ListChannelFiles(teamID, channelID string) ([]TeamFile, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	// First get the filesFolder drive item
	path := fmt.Sprintf("/teams/%s/channels/%s/filesFolder", teamID, channelID)
	body, err := s.client.Get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("get channel files folder: %w", err)
	}

	var folder struct {
		ID            string `json:"id"`
		ParentRef     struct {
			DriveID string `json:"driveId"`
		} `json:"parentReference"`
	}
	if err := json.Unmarshal(body, &folder); err != nil {
		return nil, err
	}

	// Now list the children of that folder
	childPath := fmt.Sprintf("/drives/%s/items/%s/children", folder.ParentRef.DriveID, folder.ID)
	childBody, err := s.client.Get(childPath, map[string]string{
		"$select": "id,name,webUrl,size,lastModifiedDateTime,lastModifiedBy,folder",
		"$top":    "100",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			WebURL       string `json:"webUrl"`
			Size         int64  `json:"size"`
			LastModified string `json:"lastModifiedDateTime"`
			LastModBy    struct {
				User struct {
					DisplayName string `json:"displayName"`
				} `json:"user"`
			} `json:"lastModifiedBy"`
			Folder *struct {
				ChildCount int `json:"childCount"`
			} `json:"folder"`
		} `json:"value"`
	}
	if err := json.Unmarshal(childBody, &resp); err != nil {
		return nil, err
	}

	files := make([]TeamFile, 0, len(resp.Value))
	for _, v := range resp.Value {
		files = append(files, TeamFile{
			ID:             v.ID,
			Name:           v.Name,
			WebURL:         v.WebURL,
			Size:           v.Size,
			LastModified:   v.LastModified,
			LastModifiedBy: v.LastModBy.User.DisplayName,
			IsFolder:       v.Folder != nil,
		})
	}

	slog.Info("Teams channel files listed", "team", teamID, "channel", channelID, "count", len(files))
	return files, nil
}

// GetTeamDrive returns the default drive (document library) for a team.
func (s *TeamsService) GetTeamDrive(teamID string) (string, error) {
	if !s.client.IsConfigured() {
		return "", fmt.Errorf("Graph client not configured")
	}

	path := fmt.Sprintf("/groups/%s/drive", teamID)
	body, err := s.client.Get(path, map[string]string{
		"$select": "id",
	})
	if err != nil {
		return "", err
	}

	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}
