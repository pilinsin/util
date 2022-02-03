package util

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
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
	return videoToHLS(input, output, fmpgBinPath, fprbBinPath, false, false)
}
func VideoToSilentHLS(input, output, fmpgBinPath, fprbBinPath string) (string, error){
	return videoToHLS(input, output, fmpgBinPath, fprbBinPath, true, false)
}
func VideoToAudioHLS(input, output, fmpgBinPath, fprbBinPath string) (string, error){
	return videoToHLS(input, output, fmpgBinPath, fprbBinPath, false, true)
}
func videoToHLS(input, output, fmpgBinPath, fprbBinPath string, skipA, skipV bool) (string, error){
	if output == ""{output = input}
	outExt := filepath.Ext(output)
	outNoExt := strings.TrimRight(output, outExt)
	if skipV{
		output = outNoExt + ".m4a"
	}else{
		output = outNoExt + ".m3u8"
	}

	opts := ffmpeg.Options{
		OutputFormat: StrPtr("hls"),
		VideoCodec: StrPtr("copy"),
		AudioCodec: StrPtr("copy"),
		HlsSegmentDuration: IntPtr(10),
		HlsListSize: IntPtr(0),
		HlsSegmentFilename: StrPtr(outNoExt + "_%3d.ts"),
	}
	if skipA{
		opts.SkipAudio = &skipA
	}
	if skipV{
		opts.SkipVideo = &skipV
		opts.HlsSegmentFilename = StrPtr(outNoExt + "_%3d.m4a")
	}
	fmt.Println(opts.GetStrArguments())
	conf := &ffmpeg.Config{
		FfmpegBinPath: fmpgBinPath,
		FfprobeBinPath: fprbBinPath,
		ProgressEnabled: true,
	}
	progress, err := ffmpeg.New(conf).Input(input).Output(output).WithOptions(opts).Start(opts)
	for msg := range progress{
		fmt.Println(msg)
	}
	return output, err
}
