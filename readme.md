## Kick VOD Downloader (WIP)

This is a tool to download the VOD from kick.com

---

_TODO list:_

- [x] Get the m3u8 file from kick
- [x] Parse M3U8 Master Playlist
- [x] Parse M3U8 Video Playlist
- [x] Download the video segments
- [x] Merge the video segments
- [ ] Convert the video to mp4
---

## How to use

1. Build the project
2. Run the project with the kick url as the first argument

__Or__

3. Enter the kick url like https://www.kick.com/video/d5843b1c-70d8-426c-bdaf-d69e6d90b80c in the console


---

## How to build

> go build -o kick-vod-downloader.exe