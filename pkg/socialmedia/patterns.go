package socialmedia

import "regexp"

var PlatformPatterns = map[string]*regexp.Regexp{
	"tiktok": regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|vm\.|vt\.|m\.)?(?:tiktok\.com|tiktokcdn\.com)(?:\/.*)?`),
	// "capcut":      regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?(?:capcut\.com|capcutpro\.com)(?:\/.*)?`),
	// "xiaohongshu": regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?(?:xiaohongshu\.com|xhslink\.com|xhs\.cn)(?:\/.*)?`),
	// "threads":     regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?threads\.net(?:\/.*)?`),
	// "soundcloud":  regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?(?:soundcloud\.com|snd\.sc)(?:\/.*)?`),
	// "spotify":     regexp.MustCompile(`(?:https?:\/\/)?(?:open\.)?spotify\.com(?:\/.*)?`),
	// "facebook":    regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?facebook\.com(?:\/.*)?`),
	// "instagram":   regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?instagram\.com(?:\/.*)?`),
	// "terabox":     regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?(?:terabox\.com|teraboxapp\.com)(?:\/.*)?`),
	// "snackvideo":  regexp.MustCompile(`(?:https?:\/\/)?(?:www\.|m\.)?snackvideo\.com(?:\/.*)?`),
}
