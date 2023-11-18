import { type VideoMetadata } from "../types/videoMetadata";


export async function getVideoMetadata(skip: number, firstLoad: boolean = false) {

    let videoMetadata: VideoMetadata = null;
    let url = `http://localhost:3000/metadata?skip=${skip}&firstLoad=${firstLoad}`
    let response = await fetch(url, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    });

    if (response.status == 200) {
        videoMetadata = await response.json();
    }

    return videoMetadata

}