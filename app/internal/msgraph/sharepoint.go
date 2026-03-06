package msgraph

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

// SharePointService reads SharePoint sites, lists, and documents via Microsoft Graph.
type SharePointService struct {
	client *Client
}

// NewSharePointService creates a SharePoint service backed by the given Graph client.
func NewSharePointService(client *Client) *SharePointService {
	return &SharePointService{client: client}
}

// Site represents a SharePoint site.
type Site struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	WebURL      string `json:"webUrl"`
	Description string `json:"description"`
}

// SPList represents a SharePoint list.
type SPList struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	ItemCount   int    `json:"list>contentTypesEnabled"` // nested under "list" object
	WebURL      string `json:"webUrl"`
}

// SPListItem represents a SharePoint list item with dynamic fields.
type SPListItem struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// DriveItem represents a file or folder in SharePoint document library.
type DriveItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	WebURL           string `json:"webUrl"`
	Size             int64  `json:"size"`
	IsFolder         bool
	LastModified     string `json:"lastModifiedDateTime"`
	LastModifiedBy   string
	ChildCount       int
}

// GetSiteByURL resolves a SharePoint site URL to a site ID.
// siteURL: e.g. "https://yourtenant.sharepoint.com/sites/TBS"
func (s *SharePointService) GetSiteByURL(siteURL string) (*Site, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	// Parse the URL to extract hostname and server-relative path
	u, err := url.Parse(siteURL)
	if err != nil {
		return nil, fmt.Errorf("invalid site URL: %w", err)
	}

	hostname := u.Hostname()
	serverRelativePath := strings.TrimRight(u.Path, "/")
	if serverRelativePath == "" {
		serverRelativePath = "/"
	}

	path := fmt.Sprintf("/sites/%s:%s", hostname, serverRelativePath)
	body, err := s.client.Get(path, map[string]string{
		"$select": "id,displayName,webUrl,description",
	})
	if err != nil {
		return nil, fmt.Errorf("resolve site: %w", err)
	}

	var site Site
	if err := json.Unmarshal(body, &site); err != nil {
		return nil, err
	}
	return &site, nil
}

// ListLists returns all lists in a SharePoint site.
func (s *SharePointService) ListLists(siteID string) ([]SPList, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	path := fmt.Sprintf("/sites/%s/lists", siteID)
	body, err := s.client.Get(path, map[string]string{
		"$select": "id,displayName,description,webUrl",
		"$top":    "50",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value []SPList `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}

// GetListItems reads all items from a SharePoint list with their fields.
func (s *SharePointService) GetListItems(siteID, listID string) ([]SPListItem, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	path := fmt.Sprintf("/sites/%s/lists/%s/items", siteID, listID)
	params := map[string]string{
		"$expand": "fields",
		"$top":    "500",
	}

	var allItems []SPListItem

	for path != "" {
		body, err := s.client.Get(path, params)
		if err != nil {
			return nil, err
		}

		var resp struct {
			Value    []SPListItem `json:"value"`
			NextLink string       `json:"@odata.nextLink"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		allItems = append(allItems, resp.Value...)

		// Handle pagination
		if resp.NextLink != "" {
			// NextLink is a full URL; extract the path after /v1.0
			parts := strings.SplitN(resp.NextLink, "/v1.0", 2)
			if len(parts) == 2 {
				path = parts[1]
				params = nil // params are embedded in NextLink
			} else {
				break
			}
		} else {
			break
		}
	}

	slog.Info("SharePoint list items loaded", "siteID", siteID, "listID", listID, "count", len(allItems))
	return allItems, nil
}

// ListDriveItems returns files and folders in a SharePoint document library.
func (s *SharePointService) ListDriveItems(siteID, driveID, folderPath string) ([]DriveItem, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	var path string
	if folderPath == "" || folderPath == "/" {
		path = fmt.Sprintf("/sites/%s/drives/%s/root/children", siteID, driveID)
	} else {
		path = fmt.Sprintf("/sites/%s/drives/%s/root:/%s:/children", siteID, driveID, folderPath)
	}

	body, err := s.client.Get(path, map[string]string{
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
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	items := make([]DriveItem, 0, len(resp.Value))
	for _, v := range resp.Value {
		item := DriveItem{
			ID:             v.ID,
			Name:           v.Name,
			WebURL:         v.WebURL,
			Size:           v.Size,
			LastModified:   v.LastModified,
			LastModifiedBy: v.LastModBy.User.DisplayName,
			IsFolder:       v.Folder != nil,
		}
		if v.Folder != nil {
			item.ChildCount = v.Folder.ChildCount
		}
		items = append(items, item)
	}

	return items, nil
}

// GetDrives lists document libraries (drives) for a SharePoint site.
func (s *SharePointService) GetDrives(siteID string) ([]struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WebURL      string `json:"webUrl"`
}, error) {
	if !s.client.IsConfigured() {
		return nil, fmt.Errorf("Graph client not configured")
	}

	path := fmt.Sprintf("/sites/%s/drives", siteID)
	body, err := s.client.Get(path, map[string]string{
		"$select": "id,name,description,webUrl",
	})
	if err != nil {
		return nil, err
	}

	var resp struct {
		Value []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			WebURL      string `json:"webUrl"`
		} `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Value, nil
}
