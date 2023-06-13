package autoload

/*
	You can just read the .env file on import just by doing

		import _ "github.com/titti-008/godotenv/autload"

	And bob's your mother's brother.
*/

import "github.com/titti-008/godotenv"

func init() {
	godotenv.Load()
}
