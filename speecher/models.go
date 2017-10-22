package speecher

import "github.com/FireGM/speechkit"
import "os"

var Client *speechkit.Client

func init() {
	Client = speechkit.DefaultClient(os.Getenv("SPEECHKIT_APIKEY"))
}
