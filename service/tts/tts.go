package tts

import "GoNexus/common/tts"

// BaiduTTSService 百度云TTS服务层结构体
type BaiduTTSService struct {
	TTSService *tts.BaiduTTS
}

// NewBaiduTTSService 百度云TTS服务层实例化方法
func NewBaiduTTSService() *BaiduTTSService {
	return &BaiduTTSService{
		TTSService: tts.NewBaiduTTS(),
	}
}
