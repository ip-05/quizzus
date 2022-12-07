package server

func Init() {
	// config := config.GetConfig()
	r := NewRouter()

	r.Run("0.0.0.0:3000")
}
