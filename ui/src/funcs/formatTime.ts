export function formatTime(time: number): string {
    const hours: number = Math.floor(time / 3600);
    const minutes: number = Math.floor((time - hours * 3600) / 60);
    const remainingSeconds: number = time - hours * 3600 - minutes * 60;
    const seconds: number = Math.floor(remainingSeconds);

    let formattedHours: string = hours < 10 ? `0${hours}` : `${hours}`;
    let formattedMinutes: string =
        minutes < 10 ? `0${minutes}` : `${minutes}`;
    let formattedSeconds: string =
        seconds < 10 ? `0${seconds}` : `${seconds}`;


    if (formattedHours == "NaN") {
        formattedHours = "00";
    }
    if (formattedMinutes == "NaN") {
        formattedMinutes = "00";
    }
    if (formattedSeconds == "NaN") {
        formattedSeconds = "00";
    }

    return `${formattedHours}:${formattedMinutes}:${formattedSeconds}`;
}
