package kovaaks

type Client struct {
	kovaaksPath string
	Playlists   []string
}

// New returns a new kovaaks client with the given kovaaks path
func New() (*Client, error) {
	newClient := &Client{}

	err := newClient.populateKovaaksPath()
	if err != nil {
		return nil, err
	}

	return newClient, nil
}

func stripJSONFileExtension(fileName string) string {
	return fileName[:len(fileName)-5]
}
