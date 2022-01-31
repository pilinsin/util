package util

import (
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"path/filepath"
	"strings"
	ffmpeg "github.com/floostack/transcoder/ffmpeg"
)

func LoadImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	return img, err
}
func LoadImageConfig(r io.Reader) (image.Config, error) {
	cfg, _, err := image.DecodeConfig(r)
	return cfg, err
}

func VideoToHLS(input, output, fmpgBinPath, fprbBinPath string) (string, error){
	if output == ""{output = input}
	outExt := strings.TrimLeft(filepath.Ext(output), ".")
	if outExt != "m3u8"{
		output = strings.TrimRight(output, outExt)
		output += "m3u8"
	}

	opts := ffmpeg.Options{
		OutputFormat: StrPtr("hls"),
		VideoCodec: StrPtr("copy"),
		AudioCodec: StrPtr("copy"),
		HlsSegmentDuration: IntPtr(10),
		HlsListSize: IntPtr(0),
	}
	conf := &ffmpeg.Config{
		FfmpegBinPath: fmpgBinPath,
		FfprobeBinPath: fprbBinPath,
		ProgressEnabled: true,
	}
	_, err := ffmpeg.New(conf).Input(input).Output(output).WithOptions(opts).Start(opts)
	return output, err
}