package runner

import (
	"encoding/binary"
	"log/slog"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/stealth"
	rodstream "github.com/navicstein/rod-stream"
	"gmvr.pw/boombox-web-runner/pkg/model"
	"layeh.com/gopus"
)

type RunnerService struct {
	runnerRepository RunnerRepository
	moduleRepository ModuleRepository

	encoder  *gopus.Encoder
	decoder  *gopus.Decoder
	launcher string

	logger *slog.Logger
}

func NewRunnerService(logger *slog.Logger) (*RunnerService, error) {
	var err error

	s := RunnerService{logger: logger}

	s.launcher = rodstream.MustPrepareLauncher(rodstream.LauncherArgs{}).
		NoSandbox(true).
		MustLaunch()

	s.encoder, err = gopus.NewEncoder(SampleRate, Channels, gopus.Audio)
	if err != nil {
		return nil, err
	}

	s.decoder, err = gopus.NewDecoder(SampleRate, Channels)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *RunnerService) Init(
	runnerRepository RunnerRepository,
	moduleRepository ModuleRepository,
) error {
	s.runnerRepository = runnerRepository
	s.moduleRepository = moduleRepository

	return nil
}

func (s *RunnerService) Run(runner *model.Runner, out chan<- []byte) error {
	var err error

	module, err := s.moduleRepository.GetModuleByUrl(runner.Url)
	if err != nil {
		return err
	}

	browser := rod.New().ControlURL(s.launcher).DefaultDevice(devices.Pixel2).MustConnect()

	rodstream.GrantPermissions([]string{runner.Url}, browser)

	// base64 encoded webm data
	raw := make(chan string, 2)
	opus := make(chan []byte, 2)
	s.pipe(raw, opus)

	authorize(browser, module)

	target := rodstream.MustCreatePage(browser)
	page := stealth.MustPage(browser)
	page.MustNavigate(runner.Url).WaitLoad()

	if err := rodstream.MustGetStream(target, rodstream.StreamConstraints{
		Audio:     true,
		Video:     false,
		MimeType:  "audio/webm;codecs=opus",
		FrameSize: 1000,
	}, raw); err != nil {
		s.logger.Error("cannot get stream", "error", err)
		return err
	}

	stop := make(chan bool)
	go func() {
		defer func() {
			close(out)
			s.runnerRepository.DeleteById(runner.ID)
		}()

		timecode := runner.Timecode
		play(timecode, page, module)

		tick := 0
		for {
			select {
			case <-stop:
				return
			case o := <-opus:
				if tick = (tick + 1) % 50; tick == 0 {
					if timecode, err = now(page, module); err != nil {
						return
					}
					if finished, _ := finished(page, module); finished {
						rodstream.MustStopStream(target)
						page.Close()
						browser.Close()

						tick := time.Tick(time.Second)
						ticked := 0
						for {
							select {
							case <-stop:
								return
							case <-tick:
								if ticked++; ticked == 50 {
									return
								}
								out <- []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
							}
						}
					}
				}

				data := make([]byte, 8)
				binary.BigEndian.PutUint64(data, timecode)
				out <- append(data, o...)
			}
		}
	}()

	runner.Stop = stop
	err = s.runnerRepository.Create(runner)
	if err != nil {
		s.logger.Error("cannot create runner", "error", err)
		return err
	}

	return nil
}

func (s *RunnerService) Stop(runnerId string) (*model.Runner, error) {
	runner, err := s.runnerRepository.DeleteById(runnerId)
	if err != nil {
		return nil, err
	}

	runner.Stop <- true
	return runner, nil
}
