Local test pictures folder

Put your private/local test images here for manual runs, for example:

pic2video-windows-amd64.exe render --input ./tests/fixtures/local-images --profile fhd

Files in this folder are ignored by git, except .gitkeep.

For EXIF overlay testing, include at least two images with camera metadata fields
(model, focal length, exposure time, aperture, ISO, capture date), then run:

pic2video-windows-amd64.exe render --input ./tests/fixtures/local-images --profile fhd --exif-overlay --exif-font-size 42
