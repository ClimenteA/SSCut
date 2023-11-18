import { type WebResponse } from "./webResponse"

export interface FrameInfo {
    Frame: number
    Keep: boolean
    Url: string
}

export interface VideoMetadata extends WebResponse {
    TotalSeconds: number
    Frames: FrameInfo[]
    LastSkip: number
    VideoSrc: string
}
