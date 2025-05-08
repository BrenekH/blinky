package apiunstable

import "os"

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// func removeIfExists(name string) error {
// 	err := os.Remove(name)
// 	if os.IsNotExist(err) {
// 		return nil
// 	}
// 	return err
// }
