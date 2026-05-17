package sound

import (
	"fmt"
	"log"

	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	audioWorks bool
	sounds     = make(map[int][]byte)
	devID      sdl.AudioDeviceID
	streams    []*sdl.AudioStream
)

func LoadSnd(id int, fname string) error {
	if !audioWorks {
		log.Printf("Not loading %s, as audio does not work", fname)
		return nil
	}
	_, ok := sounds[id]
	if ok {
		return fmt.Errorf("audio id %d already in use", id)
	}

	var spec sdl.AudioSpec
	wav, err := sdl.LoadWAV(fname, &spec)
	if err != nil {
		return err
	}
	sounds[id] = wav
	log.Printf("%s -> %d loaded: %+v", fname, id, spec)
	return nil
}

func Init() error {
	log.Println("Sound init...")
	driver := sdl.GetCurrentAudioDriver()
	log.Println("Current audio driver:", driver)

	devs, err := sdl.GetAudioPlaybackDevices()
	if err != nil {
		log.Println("failed to get playback devices")
		return err
	}
	if len(devs) == 0 {
		log.Println("no playback dev detected")
		return fmt.Errorf("no playback dev detected")
	}
	log.Println("playback devs", devs)

	for _, d := range devs {
		name, _ := d.Name()
		log.Println("Trying to open device:", name)
		devID, err = d.OpenAudioDevice(nil)
		if err != nil {
			log.Println("Failed to open audio device:", err)
			continue
		}
		break
	}

	err = devID.Resume()
	if err != nil {
		log.Println("Failed to start device")
		return err
	}

	log.Println("Audio init successful")
	audioWorks = true
	return nil
}

func Close() {
	audioWorks = false
	sounds = nil

	for _, stream := range streams {
		_ = stream.Clear()
		stream.Destroy()
	}
	streams = nil

	_ = devID.Pause()
	devID.Close()
	log.Println("Audio closed")
}

func PlaySnd(id int) error {
	if !audioWorks {
		log.Printf("Not playing sound %d, as audio does not work", id)
		return nil
	}
	if devID.Paused() {
		err := devID.Resume()
		log.Println("Failed to resume device:", err)
		return err
	}

	data, ok := sounds[id]
	if !ok {
		return fmt.Errorf("invalid audio id: %d", id)
	}

	var stream *sdl.AudioStream
	for _, s := range streams {
		l, err := s.Queued()
		if err != nil {
			log.Println("Failed to query stream:", err)
			continue
		}
		if l < 100 {
			// stream about to finish
			stream = s
			break
		}
	}
	if stream == nil {
		var err error
		stream, err = sdl.CreateAudioStream(nil, nil) // TODO: currently this can grow indefinitely!!!! (maybe not an issue in such a short lived application for now)
		if err != nil {
			log.Println("failed to create audio stream:", err)
			return err
		}

		err = devID.BindAudioStream(stream)
		if err != nil {
			log.Println("failed to bind audio stream:", err)
			return err
		}

		streams = append(streams, stream)
		log.Println("New audio stream created! total streams:", len(streams))
	}

	err := stream.PutData(data)
	if err != nil {
		log.Println("Failed to put data onto stream:", err)
		return err
	}
	return nil
}
