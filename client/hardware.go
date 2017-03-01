// +build !linux,!arm

package main

func readyHardware() error {
	log.Info("Hello I'm not ARM!")
	return nil
}

func closeHardware() error {
	return nil
}

func vendItem(location uint8) {
	log.Infof("==== I AM VENDING ITEM #%d! ====", location)
}
