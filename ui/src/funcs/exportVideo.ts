export async function exportVideo(outputVolume: number) {

    let options = {
        OutputVolume: outputVolume
    }

    const response = await fetch("http://localhost:3000/export-video", {
        method: "POST",
        body: JSON.stringify(options),
        headers: {
            "Content-Type": "application/json",
        },
    });

    const data = await response.json();

    if (data.Status != "success") {
        throw new Error("Failed to export video");
    } else {
        console.log("video is exported in the background")
    }
}