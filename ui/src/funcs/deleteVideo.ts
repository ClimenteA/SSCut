

export async function deleteVideo() {
    const response = await fetch("http://localhost:3000/delete-video", {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json",
        },
    });

    const data = await response.json();

    if (data.Status != "success") {
        throw new Error("Failed to delete video");
    }
}