import { type TotalSecondsWebResponse } from "../types/webResponse";

export async function getVideoTotalSeconds(): Promise<number> {
    const response = await fetch("http://localhost:3000/total-seconds", {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    });

    const data: TotalSecondsWebResponse = await response.json();

    if (data.Status == "success") {
        return data.TotalSeconds;
    } else {
        throw new Error("Failed to fetch total seconds");
    }
}
