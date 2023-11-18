import { type FrameInfo } from "../types/videoMetadata";


export async function setVideoMetadata(frame: FrameInfo) {

    frame.Keep = !frame.Keep

    const response = await fetch("http://localhost:3000/save-cutted", {
        method: "POST",
        body: JSON.stringify(frame),
        headers: {
            "Content-Type": "application/json",
        },
    });

    const data = await response.json();

    if (data.Status != "success") {
        console.log(data);
        throw new Error("Failed to set frame state");
    }

}
