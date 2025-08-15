package db

func AddAuthData(username string, password string) error {
	return nil
}

func CheckAuthData(username string, password string) (bool, error) {
	return true, nil
}

func AddAudioTask(username string, audio string) error {
	return nil
}

func GetStatusTask(username string) (string, error) {
	return "", nil
}

func GetResultTask(username string) (string, error) {
	return "", nil
}
