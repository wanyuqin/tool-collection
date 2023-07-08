package configs

import (
	"github.com/wanyuqin/tool-collection/logger"
	"log"
	"testing"
)

func TestSaveDownloadSettings(t *testing.T) {
	logger.InitLogger()
	err := SaveDownloadSettings(DownloadConfig{
		Path: "/Users/ethanleo",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}
