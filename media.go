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
	"strconv"
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


func VideoToHLS(input, output, fmpgBinPath, fprbBinPath string, additionalArgs ...string) (string, error){
	if output == ""{output = input}
	outExt := filepath.Ext(output)
	outNoExt := strings.TrimRight(output, outExt)
	output = outNoExt + ".m3u8"
	
	tr := NewFfmpeg(fmpgBinPath, fprbBinPath)
	args := []string{
		"-c:v copy",
		"-c:a copy",
		"-hls_time 10",
		"-hls_list_size 0",
		"-hls_segment_filename " + outNoExt + "_%3d.ts",
		"-f hls",
	}
	args = append(args, additionalArgs...)
	return tr.Transcode(input, output, args...)
}

type ffTranscoder struct{
	cfg *ffmpeg.Config
}
func NewFfmpeg(ffmpegBinPath, ffProbeBinPath string) *ffTranscoder{
	return &ffTranscoder{
		cfg: &ffmpeg.Config{
			FfmpegBinPath: ffmpegBinPath,
			FfprobeBinPath: ffProbeBinPath,
			ProgressEnabled: true,
		},
	}
}
func (self ffTranscoder) Transcode(input, output string, args ...string) (string, error){
	if output == ""{output = input}

	opts := ffOpt.newOptions(args...)
	fmt.Println(opts.GetStrArguments())

	progress, err := ffmpeg.New(self.cfg).Input(input).Output(output).WithOptions(opts).Start(opts)
	for msg := range progress{
		fmt.Println(msg)
	}
	return output, err
}

func (self ffmpegOption) newOptions(args ...string) *ffmpeg.Options{
	opts := &ffmpeg.Options{}
	for _, arg := range args{
		kv := strings.Fields(arg)
		if optFunc, ok := ffmpegOptionsMap[kv[0]]; ok{
			if len(kv) == 1{
				optFunc("")(opts)
			}else{
				v := strings.Join(kv[1:], " ")
				optFunc(v)(opts)
			}
		}
	}
	return opts
}

var ffmpegOptionsMap = map[string] func(string) ffmpegSettings{
	"-aspect": ffOpt.aspect,
	"-s": ffOpt.resolution,
	"-b:v": ffOpt.videoBitRate,
	"-bt": ffOpt.videoBitRateTolerance,
	"-maxrate": ffOpt.videoMaxBitRate,
	"-minrate": ffOpt.videoMaxBitRate,
	"-c:v": ffOpt.videoCodec,
	"-vframes": ffOpt.vframes,
	"-r": ffOpt.frameRate,
	"-ar": ffOpt.audioRate,
	"-g": ffOpt.keyframeInterval,
	"-c:a": ffOpt.audioCodec,
	"-ab": ffOpt.audioBitRate,
	"-ac": ffOpt.audioChannels,
	"-q:a": ffOpt.audioVariableBitRate,
	"-bufsize": ffOpt.bufferSize,
	"-threads": ffOpt.threads,
	"-preset": ffOpt.preset,
	"-tune": ffOpt.tune,
	"-profile:a": ffOpt.audioProfile,
	"-profile:v": ffOpt.videoProfile,
	"-target": ffOpt.target,
	"-t": ffOpt.duration,
	"-qscale": ffOpt.qscale,
	"-crf": ffOpt.crf,
	"-strict": ffOpt.strict,
	"-muxdelay": ffOpt.muxDelay,
	"-ss": ffOpt.seekTime,
	"-seek_timestamp": ffOpt.seekUsingTimestamp,
	"-movflags": ffOpt.movFlags,
	"-hide_banner": ffOpt.hideBanner,
	"-f": ffOpt.outputFormat,
	"-copyts": ffOpt.copyTs,
	"-re": ffOpt.nativeFrameRateInput,
	"-itsoffset": ffOpt.inputInitialOffset,
	"-rtmp_live": ffOpt.rtmpLive,
	"-hls_playlist_type": ffOpt.hlsPlaylistType,
	"-hls_list_size": ffOpt.hlsListSize,
	"-hls_time": ffOpt.hlsSegmentDuration,
	"-master_pl_name": ffOpt.hlsMasterPlaylistName,
	"-hls_segment_filename": ffOpt.hlsSegmentFilename,
	"-method": ffOpt.httpMethod,
	"-multiple_requests": ffOpt.httpKeepAlive,
	"-hwaccel": ffOpt.hwaccel,
	"-vf": ffOpt.videoFilter,
	"-af": ffOpt.audioFilter,
	"-vn": ffOpt.skipVideo,
	"-an": ffOpt.skipAudio,
	"-compression_level": ffOpt.compressionLevel,
	"-map_metadata": ffOpt.mapMetadata,
	"-hls_key_info_file": ffOpt.encryptionKey,
	"-bf": ffOpt.bframe,
	"-pix_fmt": ffOpt.pixFmt,
	"-y": ffOpt.overwrite,
}
type ffmpegSettings func(*ffmpeg.Options)
type ffmpegOption struct{}
var ffOpt ffmpegOption
func parseIntOption(prm string) *int{
	if n, err := strconv.Atoi(prm); err != nil{
		return nil
	}else{
		return &n
	}
}
func parseUint32Option(prm string) *uint32{
	if n, err := strconv.ParseUint(prm, 10, 32); err != nil{
		return nil
	}else{
		u32n := uint32(n)
		return &u32n
	}
}
func (self ffmpegOption)aspect(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Aspect = &prm
	}
}
func (self ffmpegOption)resolution(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Resolution = &prm
	}
}
func (self ffmpegOption)videoBitRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoBitRate = &prm
	}
}
func (self ffmpegOption)videoBitRateTolerance(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoBitRateTolerance = parseIntOption(prm)
	}
}
func (self ffmpegOption)videoMaxBitRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoMaxBitRate = parseIntOption(prm)
	}
}
func (self ffmpegOption)videoMinBitRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoMinBitrate = parseIntOption(prm)
	}
}
func (self ffmpegOption)videoCodec(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoCodec = &prm
	}
}
func (self ffmpegOption)vframes(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Vframes = parseIntOption(prm)
	}
}
func (self ffmpegOption)frameRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.FrameRate = parseIntOption(prm)
	}
}
func (self ffmpegOption)audioRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioRate = parseIntOption(prm)
	}
}
func (self ffmpegOption)keyframeInterval(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.KeyframeInterval = parseIntOption(prm)
	}
}
func (self ffmpegOption)audioCodec(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioCodec = &prm
	}
}
func (self ffmpegOption)audioBitRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioBitrate = &prm
	}
}
func (self ffmpegOption)audioChannels(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioChannels = parseIntOption(prm)
	}
}
func (self ffmpegOption)audioVariableBitRate(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioVariableBitrate = BoolPtr(true)
	}
}
func (self ffmpegOption)bufferSize(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.BufferSize = parseIntOption(prm)
	}
}
func (self ffmpegOption)threads(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Threads = parseIntOption(prm)
	}
}
func (self ffmpegOption)preset(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Preset = &prm
	}
}
func (self ffmpegOption)tune(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Tune = &prm
	}
}
func (self ffmpegOption)audioProfile(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioProfile = &prm
	}
}
func (self ffmpegOption)videoProfile(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoProfile = &prm
	}
}
func (self ffmpegOption)target(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Target = &prm
	}
}
func (self ffmpegOption)duration(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Duration = &prm
	}
}
func (self ffmpegOption)qscale(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Qscale = parseUint32Option(prm)
	}
}
func (self ffmpegOption)crf(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Crf = parseUint32Option(prm)
	}
}
func (self ffmpegOption)strict(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Strict = parseIntOption(prm)
	}
}
func (self ffmpegOption)muxDelay(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.MuxDelay = &prm
	}
}
func (self ffmpegOption)seekTime(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.SeekTime = &prm
	}
}
func (self ffmpegOption)seekUsingTimestamp(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.SeekUsingTimestamp = BoolPtr(true)
	}
}
func (self ffmpegOption)movFlags(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.MovFlags = &prm
	}
}
func (self ffmpegOption)hideBanner(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HideBanner = BoolPtr(true)
	}
}
func (self ffmpegOption)outputFormat(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.OutputFormat = &prm
	}
}
func (self ffmpegOption)copyTs(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.CopyTs = BoolPtr(true)
	}
}
func (self ffmpegOption)nativeFrameRateInput(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.NativeFramerateInput = BoolPtr(true)
	}
}
func (self ffmpegOption)inputInitialOffset(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.InputInitialOffset = &prm
	}
}
func (self ffmpegOption)rtmpLive(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.RtmpLive = &prm
	}
}
func (self ffmpegOption)hlsPlaylistType(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HlsPlaylistType = &prm
	}
}
func (self ffmpegOption)hlsListSize(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HlsListSize = parseIntOption(prm)
	}
}
func (self ffmpegOption)hlsSegmentDuration(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HlsSegmentDuration = parseIntOption(prm)
	}
}
func (self ffmpegOption)hlsMasterPlaylistName(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HlsMasterPlaylistName = &prm
	}
}
func (self ffmpegOption)hlsSegmentFilename(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HlsSegmentFilename = &prm
	}
}
func (self ffmpegOption)httpMethod(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HTTPMethod = &prm
	}
}
func (self ffmpegOption)httpKeepAlive(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.HTTPKeepAlive = BoolPtr(true)
	}
}
func (self ffmpegOption)hwaccel(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Hwaccel = &prm
	}
}
func (self ffmpegOption)videoFilter(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.VideoFilter = &prm
	}
}
func (self ffmpegOption)audioFilter(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.AudioFilter = &prm
	}
}
func (self ffmpegOption)skipVideo(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.SkipVideo = BoolPtr(true)
	}
}
func (self ffmpegOption)skipAudio(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.SkipAudio = BoolPtr(true)
	}
}
func (self ffmpegOption)compressionLevel(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.CompressionLevel = parseIntOption(prm)
	}
}
func (self ffmpegOption)mapMetadata(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.MapMetadata = &prm
	}
}
func (self ffmpegOption)encryptionKey(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.EncryptionKey = &prm
	}
}
func (self ffmpegOption)bframe(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Bframe = parseIntOption(prm)
	}
}
func (self ffmpegOption)pixFmt(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.PixFmt = &prm
	}
}
func (self ffmpegOption)overwrite(prm string) ffmpegSettings{
	return func(opt *ffmpeg.Options){
		opt.Overwrite = BoolPtr(true)
	}
}