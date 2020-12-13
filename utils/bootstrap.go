package utils

import "oauth/config"

func StartUp(datastoreSource *config.Datasource) {

	// Start a MongoDB session
	CreateDBConnection(datastoreSource.ConnectionURL, datastoreSource.Databasename)
}
