package global

import "github.com/rs/zerolog/log"

var (
	volume                    float64 = 0.2
	onVolumeChangeSubscribers []func()
)

func Volume() float64 {
	return volume
}

func SetVolume(value float64) {
	if volume > 1 {
		volume = 1
	}
	if volume < 0 {
		volume = 0
	}
	if volume == value {
		return
	}

	volume = value
	for _, event := range onVolumeChangeSubscribers {
		event()
	}
	log.Debug().Msgf("volume is now %0.2f", volume)
}

func SubscribeOnVolumeChange(event func()) {
	onVolumeChangeSubscribers = append(onVolumeChangeSubscribers, event)
}
