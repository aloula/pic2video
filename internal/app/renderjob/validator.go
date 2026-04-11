package renderjob

import "os"

func ValidateJob(job RenderJob) error {
	if len(job.InputAssets) < 1 {
		return &ClassifiedError{Class: ErrInputValidation, Msg: "at least 1 valid media asset is required"}
	}
	if job.Profile.Width <= 0 || job.Profile.Height <= 0 {
		return &ClassifiedError{Class: ErrInvalidArguments, Msg: "invalid profile dimensions"}
	}
	if job.ImageDurationSec <= 0 {
		return &ClassifiedError{Class: ErrInvalidArguments, Msg: "image duration must be > 0"}
	}
	if job.TransitionDurationSec <= 0 || job.TransitionDurationSec >= job.ImageDurationSec {
		return &ClassifiedError{Class: ErrInvalidArguments, Msg: "transition duration must be > 0 and less than image duration"}
	}
	if job.OutputPath == "" {
		return &ClassifiedError{Class: ErrInvalidArguments, Msg: "output path is required"}
	}
	if job.OutputFPS < 24 || job.OutputFPS > 60 {
		return &ClassifiedError{Class: ErrInvalidArguments, Msg: "output fps must be between 24 and 60"}
	}
	if _, err := os.Stat(job.OutputPath); err == nil && !job.Overwrite {
		return &ClassifiedError{Class: ErrInputValidation, Msg: "output file already exists (use --overwrite)"}
	}
	return nil
}
