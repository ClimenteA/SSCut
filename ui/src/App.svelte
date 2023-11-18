<script lang="ts">
    import { onMount } from "svelte";
    import { getVideoMetadata } from "./funcs/getVideoMetadata";
    import { setVideoMetadata } from "./funcs/setVideoMetadata";
    import { type VideoMetadata } from "./types/videoMetadata";
    import UploadVideo from "./lib/UploadVideo.svelte";
    import { formatTime } from "./funcs/formatTime";
    import { deleteVideo } from "./funcs/deleteVideo";
    import { exportVideo } from "./funcs/exportVideo";

    let videoElement: HTMLVideoElement;
    let videoMetadata: VideoMetadata;
    let videoPlaying = false;
    let inputRange: HTMLInputElement;
    let remainingTime = "00:00:00";
    let selectedFrame = 1;
    let playRate = 1;
    let roundedCurrentTime = 0;
    let currentTime = 0;
    let skip = 0;
    let outputVolume = 20;
    let msgSuccess: string;
    let msgFailed: string;
    let deleteVideoClicked = false;
    let videoIsBeingExported = false;
    let videoIsBeingDeleted = false;
    let videoIsBeingCutted = false;
    let videoIsBeingSkipped = false;
    let increaseOutputVolume = false;
    let mute = false;

    function showFailedDialog(message: string) {
        msgFailed = message;
        setTimeout(() => {
            msgFailed = null;
        }, 1000);
    }

    function showSuccessDialog(message: string) {
        msgSuccess = message;
        setTimeout(() => {
            msgSuccess = null;
        }, 1000);
    }

    function updateCurrentTimeVariables() {
        currentTime = videoElement.currentTime;
        remainingTime = formatTime(videoMetadata.TotalSeconds - currentTime);
        roundedCurrentTime = Math.round(currentTime);
        if (currentTime >= videoMetadata.TotalSeconds) {
            videoPlaying = false;
            remainingTime = "00:00:00";
        }
    }

    async function nextFrames() {
        skip += 7;
        if (skip >= videoMetadata.TotalSeconds) {
            skip = videoMetadata.TotalSeconds - videoMetadata.Frames.length;
            console.log("No next frames to load.", skip);
            return;
        }
        videoMetadata = await getVideoMetadata(skip);
        console.log("Loaded next frames", skip);
    }

    async function previousFrames() {
        if (skip <= 0) {
            skip = 0;
            console.log("No previous frames to load.", skip);
            return;
        }
        skip -= 7;
        videoMetadata = await getVideoMetadata(skip);
        console.log("Loaded previous frames", skip);
    }

    function frameIsVizible() {
        let visible = false;
        for (let frame of videoMetadata.Frames) {
            if (roundedCurrentTime == frame.Frame) {
                visible = true;
                break;
            }
        }
        return visible;
    }

    async function selectFrame(frame: number) {
        selectedFrame = frame;
        console.log("selectedFrame", selectedFrame);

        let lastFrame =
            videoMetadata.Frames[videoMetadata.Frames.length - 1].Frame;
        if (selectedFrame == lastFrame + 1) {
            await nextFrames();
            return;
        }
        let firstFrame = videoMetadata.Frames[0].Frame;
        if (selectedFrame == firstFrame - 1) {
            await previousFrames();
            return;
        }

        if (!frameIsVizible()) {
            console.log("bring into view frames", selectedFrame, skip);
            skip = selectedFrame - 1;
            videoMetadata = await getVideoMetadata(skip);
        }
    }

    async function manuallySelectFrame(frame: number) {
        console.log("manuallySelectFrame:", frame);
        videoElement.currentTime = frame;
    }

    async function seekVideo() {
        videoElement.currentTime = parseFloat(inputRange.value);
    }

    function toggleMute() {
        mute = !mute;
        videoElement.muted = mute;
    }

    function toggleKeyMute(event: KeyboardEvent) {
        if (event.key === "M" || event.key === "m") {
            toggleMute();
        }
    }

    function toggleVideoPlaying() {
        videoPlaying = !videoPlaying;
        videoPlaying ? videoElement.play() : videoElement.pause();
        console.log(videoPlaying ? "Video Playing" : "Video Paused");
    }

    function toggleSpaceVideoPlaying(event: KeyboardEvent) {
        if (event.code === "Space") {
            toggleVideoPlaying();
        }
    }

    async function updateTime() {
        videoIsBeingExported = true;
        videoIsBeingDeleted = true;
        videoIsBeingCutted = true;
        videoIsBeingSkipped = true;

        updateCurrentTimeVariables();
        inputRange.value = String(currentTime);
        if (selectedFrame != roundedCurrentTime) {
            selectedFrame = roundedCurrentTime;
            await selectFrame(selectedFrame);
        }

        videoIsBeingExported = false;
        videoIsBeingDeleted = false;
        videoIsBeingCutted = false;
        videoIsBeingSkipped = false;
    }

    async function cutOrRestoreFrame() {
        for (let frame of videoMetadata.Frames) {
            if (frame.Frame == selectedFrame) {
                await setVideoMetadata(frame);
                videoMetadata = await getVideoMetadata(skip);
            }
        }
    }

    async function pressCtoCutFrame(event: KeyboardEvent) {
        if (event.key == "c" || event.key == "C") {
            await cutOrRestoreFrame();
        }
    }

    function skipFoward() {
        videoElement.currentTime = roundedCurrentTime + 1;
    }

    function skipBackward() {
        videoElement.currentTime = roundedCurrentTime - 1;
    }

    function handlePlayRateChange(event: Event) {
        const target = event.target as HTMLSelectElement;
        playRate = parseFloat(target.value);
        videoElement.playbackRate = playRate;
        console.log("playRate", playRate);
    }

    function skipViaArrows(event: KeyboardEvent) {
        if (event.key === "ArrowRight") {
            skipFoward();
        }

        if (event.key === "ArrowLeft") {
            skipBackward();
        }
    }

    async function skipViewedFrames(event: KeyboardEvent) {
        if (event.key === ",") {
            await previousFrames();
        }

        if (event.key === ".") {
            await nextFrames();
        }
    }

    async function handleKeys(event: KeyboardEvent) {
        toggleSpaceVideoPlaying(event);
        await pressCtoCutFrame(event);
        skipViaArrows(event);
        await skipViewedFrames(event);
        toggleKeyMute(event);
    }

    onMount(async () => {
        videoMetadata = await getVideoMetadata(skip, true);
        console.log("Mounted: ", videoMetadata);
        remainingTime = formatTime(videoMetadata.TotalSeconds);
        skip = videoMetadata.LastSkip;
        selectFrame(skip + 1);
    });
</script>

<svelte:window on:keydown={handleKeys} />

<main class="container">
    {#if videoMetadata === null}
        <UploadVideo />
    {:else if videoMetadata !== null && videoMetadata}
        <section id="video" class="d-flex justify-content-around mt-4">
            <video
                class="video border-muted shadow-lg"
                preload="auto"
                on:timeupdate={updateTime}
                on:click={toggleVideoPlaying}
                bind:this={videoElement}
            >
                <source src={videoMetadata.VideoSrc} type="video/mp4" />
                <track kind="captions" />
            </video>
        </section>

        <section id="progress" class="mt-4">
            <input
                id="video-progress-bar"
                class="form-range"
                type="range"
                min="0"
                max={videoMetadata.TotalSeconds}
                value={currentTime}
                bind:this={inputRange}
                on:input={seekVideo}
            />

            <div class="d-flex gap-1 justify-content-center mb-4">
                <span id="remaining-time" class="badge bg-gradient fs-5"
                    >{remainingTime}</span
                >
            </div>
        </section>

        <section
            id="frames"
            class="d-flex gap-2 border border-dark rounded-2 align-items-center justify-content-center w-100 px-2 overflow-auto bg-gradient h-150px
         "
        >
            {#each videoMetadata.Frames as frame (frame.Frame)}
                <div class="d-flex flex-column align-items-center">
                    <small class="text-white">{formatTime(frame.Frame)}</small>
                    <button
                        id={"frame-btn-" + String(frame.Frame)}
                        class="btn p-0"
                        class:selected={selectedFrame === frame.Frame}
                        on:click={() => manuallySelectFrame(frame.Frame)}
                    >
                        <img
                            class="default-img h-100px"
                            src={frame.Url}
                            alt={"frame-" + String(frame.Frame)}
                        />
                    </button>
                </div>
            {/each}
        </section>

        <section id="menu" class="d-flex gap-2 justify-content-between">
            <button
                disabled={videoIsBeingSkipped}
                id="previous"
                title="Press <"
                class="btn btn-outline-light btn-sm"
                on:click={previousFrames}
            >
                {#if videoIsBeingSkipped}
                    <span
                        class="spinner-border spinner-border-sm"
                        role="status"
                        aria-hidden="true"
                    />
                {:else}
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        viewBox="0 0 16 16"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M11.354 1.646a.5.5 0 0 1 0 .708L5.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0z"
                        />
                    </svg>
                {/if}
            </button>

            <section id="controls" class="d-flex gap-2 mt-4">
                {#if !videoPlaying}
                    <button
                        title="Play (Press Space)"
                        class="btn btn-primary btn-lg"
                        on:click={toggleVideoPlaying}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="m11.596 8.697-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393z"
                            />
                        </svg>
                    </button>
                {:else if videoPlaying}
                    <button
                        title="Pause (Press Space)"
                        class="btn btn-primary btn-lg"
                        on:click={toggleVideoPlaying}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M5.5 3.5A1.5 1.5 0 0 1 7 5v6a1.5 1.5 0 0 1-3 0V5a1.5 1.5 0 0 1 1.5-1.5zm5 0A1.5 1.5 0 0 1 12 5v6a1.5 1.5 0 0 1-3 0V5a1.5 1.5 0 0 1 1.5-1.5z"
                            />
                        </svg>
                    </button>
                {/if}

                {#if mute}
                    <button
                        title="Unmute (Press M)"
                        class="btn btn-primary btn-lg"
                        on:click={toggleMute}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M6.717 3.55A.5.5 0 0 1 7 4v8a.5.5 0 0 1-.812.39L3.825 10.5H1.5A.5.5 0 0 1 1 10V6a.5.5 0 0 1 .5-.5h2.325l2.363-1.89a.5.5 0 0 1 .529-.06zm7.137 2.096a.5.5 0 0 1 0 .708L12.207 8l1.647 1.646a.5.5 0 0 1-.708.708L11.5 8.707l-1.646 1.647a.5.5 0 0 1-.708-.708L10.793 8 9.146 6.354a.5.5 0 1 1 .708-.708L11.5 7.293l1.646-1.647a.5.5 0 0 1 .708 0z"
                            />
                        </svg>
                    </button>
                {:else}
                    <button
                        title="Mute (Press M)"
                        class="btn btn-primary btn-lg"
                        on:click={toggleMute}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M11.536 14.01A8.473 8.473 0 0 0 14.026 8a8.473 8.473 0 0 0-2.49-6.01l-.708.707A7.476 7.476 0 0 1 13.025 8c0 2.071-.84 3.946-2.197 5.303l.708.707z"
                            />
                            <path
                                d="M10.121 12.596A6.48 6.48 0 0 0 12.025 8a6.48 6.48 0 0 0-1.904-4.596l-.707.707A5.483 5.483 0 0 1 11.025 8a5.483 5.483 0 0 1-1.61 3.89l.706.706z"
                            />
                            <path
                                d="M8.707 11.182A4.486 4.486 0 0 0 10.025 8a4.486 4.486 0 0 0-1.318-3.182L8 5.525A3.489 3.489 0 0 1 9.025 8 3.49 3.49 0 0 1 8 10.475l.707.707zM6.717 3.55A.5.5 0 0 1 7 4v8a.5.5 0 0 1-.812.39L3.825 10.5H1.5A.5.5 0 0 1 1 10V6a.5.5 0 0 1 .5-.5h2.325l2.363-1.89a.5.5 0 0 1 .529-.06z"
                            />
                        </svg>
                    </button>
                {/if}

                <button
                    disabled={videoIsBeingCutted}
                    title="Cut (Press C)"
                    class="btn btn-warning btn-lg"
                    on:click={async () => {
                        videoIsBeingCutted = true;
                        await cutOrRestoreFrame();
                        videoIsBeingCutted = false;
                    }}
                >
                    {#if videoIsBeingCutted}
                        <span
                            class="spinner-border spinner-border-sm"
                            role="status"
                            aria-hidden="true"
                        />
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M3.5 3.5c-.614-.884-.074-1.962.858-2.5L8 7.226 11.642 1c.932.538 1.472 1.616.858 2.5L8.81 8.61l1.556 2.661a2.5 2.5 0 1 1-.794.637L8 9.73l-1.572 2.177a2.5 2.5 0 1 1-.794-.637L7.19 8.61 3.5 3.5zm2.5 10a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0zm7 0a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0z"
                            />
                        </svg>
                    {/if}
                </button>

                <button
                    disabled={videoIsBeingSkipped}
                    title="Skip forward (Press Arrow Right)"
                    class="btn btn-primary btn-lg"
                    on:click={skipFoward}
                >
                    {#if videoIsBeingSkipped}
                        <span
                            class="spinner-border spinner-border-sm"
                            role="status"
                            aria-hidden="true"
                        />
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M8 3a5 5 0 1 0 4.546 2.914.5.5 0 0 1 .908-.417A6 6 0 1 1 8 2v1z"
                            />
                            <path
                                d="M8 4.466V.534a.25.25 0 0 1 .41-.192l2.36 1.966c.12.1.12.284 0 .384L8.41 4.658A.25.25 0 0 1 8 4.466z"
                            />
                        </svg>
                    {/if}
                </button>

                <button
                    disabled={videoIsBeingSkipped}
                    title="Skip backward (Press Arrow Left)"
                    class="btn btn-primary btn-lg"
                    on:click={skipBackward}
                >
                    {#if videoIsBeingSkipped}
                        <span
                            class="spinner-border spinner-border-sm"
                            role="status"
                            aria-hidden="true"
                        />
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                fill-rule="evenodd"
                                d="M8 3a5 5 0 1 1-4.546 2.914.5.5 0 0 0-.908-.417A6 6 0 1 0 8 2v1z"
                            />
                            <path
                                d="M8 4.466V.534a.25.25 0 0 0-.41-.192L5.23 2.308a.25.25 0 0 0 0 .384l2.36 1.966A.25.25 0 0 0 8 4.466z"
                            />
                        </svg>
                    {/if}
                </button>

                <select
                    title="Playback speed"
                    id="playRate"
                    class="form-select"
                    on:change={handlePlayRateChange}
                >
                    <option value="1">Play at 1x</option>
                    <option value="1.5">Play at 1.5x</option>
                    <option value="2">Play at 2x</option>
                    <option value="2.5">Play at 2.5x</option>
                    <option value="3">Play at 3x</option>
                </select>
            </section>

            <button
                disabled={videoIsBeingSkipped}
                id="next"
                title="Press >"
                class="btn btn-outline-light btn-sm"
                on:click={nextFrames}
            >
                {#if videoIsBeingSkipped}
                    <span
                        class="spinner-border spinner-border-sm"
                        role="status"
                        aria-hidden="true"
                    />
                {:else}
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        viewBox="0 0 16 16"
                    >
                        <path
                            fill-rule="evenodd"
                            d="M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z"
                        />
                    </svg>
                {/if}
            </button>
        </section>

        <div class="d-flex gap-4 mt-5 justify-content-end">
            <button
                disabled={videoIsBeingExported}
                title="Delete current video"
                class="btn btn-primary btn-sm"
                data-bs-target="#export-options"
                data-bs-toggle="modal"
            >
                {#if videoIsBeingExported}
                    <span
                        class="spinner-border spinner-border-sm"
                        role="status"
                        aria-hidden="true"
                    />
                {:else}
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        viewBox="0 0 16 16"
                    >
                        <path
                            d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"
                        />
                        <path
                            d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708l3 3z"
                        />
                    </svg>
                {/if}
                Export Video
            </button>

            {#if !deleteVideoClicked}
                <button
                    disabled={videoIsBeingDeleted}
                    title="Delete current video"
                    class="btn btn-secondary btn-sm"
                    on:click={() => (deleteVideoClicked = true)}
                >
                    {#if videoIsBeingDeleted}
                        <span
                            class="spinner-border spinner-border-sm"
                            role="status"
                            aria-hidden="true"
                        />
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M6.5 1h3a.5.5 0 0 1 .5.5v1H6v-1a.5.5 0 0 1 .5-.5ZM11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3A1.5 1.5 0 0 0 5 1.5v1H2.506a.58.58 0 0 0-.01 0H1.5a.5.5 0 0 0 0 1h.538l.853 10.66A2 2 0 0 0 4.885 16h6.23a2 2 0 0 0 1.994-1.84l.853-10.66h.538a.5.5 0 0 0 0-1h-.995a.59.59 0 0 0-.01 0H11Zm1.958 1-.846 10.58a1 1 0 0 1-.997.92h-6.23a1 1 0 0 1-.997-.92L3.042 3.5h9.916Zm-7.487 1a.5.5 0 0 1 .528.47l.5 8.5a.5.5 0 0 1-.998.06L5 5.03a.5.5 0 0 1 .47-.53Zm5.058 0a.5.5 0 0 1 .47.53l-.5 8.5a.5.5 0 1 1-.998-.06l.5-8.5a.5.5 0 0 1 .528-.47ZM8 4.5a.5.5 0 0 1 .5.5v8.5a.5.5 0 0 1-1 0V5a.5.5 0 0 1 .5-.5Z"
                            />
                        </svg>
                    {/if}
                    Delete video
                </button>
            {:else}
                <div class="d-flex gap-2">
                    <button
                        class="btn btn-danger btn-sm"
                        on:click={async () => {
                            videoIsBeingDeleted = true;
                            await deleteVideo();
                            videoMetadata = await getVideoMetadata(skip);
                            remainingTime = formatTime(
                                videoMetadata.TotalSeconds
                            );
                            videoIsBeingDeleted = false;
                            deleteVideoClicked = false;
                        }}
                    >
                        {#if videoIsBeingDeleted}
                            <span
                                class="spinner-border spinner-border-sm"
                                role="status"
                                aria-hidden="true"
                            />
                        {/if}
                        Delete it already!
                    </button>
                    <button
                        class="btn btn-secondary btn-sm"
                        on:click={() => (deleteVideoClicked = false)}
                    >
                        {#if videoIsBeingDeleted}
                            <span
                                class="spinner-border spinner-border-sm"
                                role="status"
                                aria-hidden="true"
                            />
                        {/if}
                        Nup keep it.
                    </button>
                </div>
            {/if}

            <button
                class="btn btn-secondary btn-sm"
                data-bs-target="#help"
                data-bs-toggle="modal"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="currentColor"
                    viewBox="0 0 16 16"
                >
                    <path
                        fill-rule="evenodd"
                        d="M4.475 5.458c-.284 0-.514-.237-.47-.517C4.28 3.24 5.576 2 7.825 2c2.25 0 3.767 1.36 3.767 3.215 0 1.344-.665 2.288-1.79 2.973-1.1.659-1.414 1.118-1.414 2.01v.03a.5.5 0 0 1-.5.5h-.77a.5.5 0 0 1-.5-.495l-.003-.2c-.043-1.221.477-2.001 1.645-2.712 1.03-.632 1.397-1.135 1.397-2.028 0-.979-.758-1.698-1.926-1.698-1.009 0-1.71.529-1.938 1.402-.066.254-.278.461-.54.461h-.777ZM7.496 14c.622 0 1.095-.474 1.095-1.09 0-.618-.473-1.092-1.095-1.092-.606 0-1.087.474-1.087 1.091S6.89 14 7.496 14Z"
                    />
                </svg>
            </button>
        </div>
    {/if}
</main>

<div
    class="modal fade"
    id="export-options"
    tabindex="-1"
    aria-labelledby="exportModalLabel"
    aria-hidden="true"
>
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="exportModalLabel">
                    Video Export Options
                </h5>
                <button
                    type="button"
                    class="btn-close"
                    data-bs-dismiss="modal"
                    aria-label="Close"
                />
            </div>
            <div class="modal-body my-3">
                <div
                    class="d-flex gap-2 justify-content-start align-items-center"
                >
                    <label for="output-volume" class="form-label"
                        >Increase volume by dB (decibels):
                        <input
                            bind:value={outputVolume}
                            type="number"
                            disabled={!increaseOutputVolume}
                            class="form-control"
                            id="output-volume"
                        />
                    </label>
                    <label class="mt-3">
                        <input
                            class="form-check-input"
                            type="checkbox"
                            bind:checked={increaseOutputVolume}
                        />
                    </label>
                </div>
            </div>
            <div class="modal-footer">
                <button
                    disabled={videoIsBeingExported}
                    class="btn btn-secondary"
                    data-bs-dismiss="modal">Close</button
                >
                <button
                    disabled={videoIsBeingExported}
                    class="btn btn-primary"
                    on:click={async () => {
                        videoIsBeingExported = true;
                        if (!increaseOutputVolume) {
                            outputVolume = 0;
                        }
                        await exportVideo(outputVolume);
                        videoIsBeingExported = false;
                    }}
                >
                    {#if videoIsBeingExported}
                        <span
                            class="spinner-border spinner-border-sm"
                            role="status"
                            aria-hidden="true"
                        />
                    {:else}
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="currentColor"
                            viewBox="0 0 16 16"
                        >
                            <path
                                d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"
                            />
                            <path
                                d="M7.646 11.854a.5.5 0 0 0 .708 0l3-3a.5.5 0 0 0-.708-.708L8.5 10.293V1.5a.5.5 0 0 0-1 0v8.793L5.354 8.146a.5.5 0 1 0-.708.708l3 3z"
                            />
                        </svg>
                    {/if}
                    Export Video
                </button>
            </div>
        </div>
    </div>
</div>

<div
    class="modal fade"
    id="help"
    tabindex="-1"
    aria-labelledby="helpModalLabel"
    aria-hidden="true"
>
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="helpModalLabel">Help</h5>
                <button
                    type="button"
                    class="btn-close"
                    data-bs-dismiss="modal"
                    aria-label="Close"
                />
            </div>
            <div class="modal-body my-3">
                <p>Here are some useful shortcuts:</p>
                <ul>
                    <li>Space - Play/Pause</li>
                    <li>M - Mute/Unmute</li>
                    <li>C - Cut/Redo Cut</li>
                    <li>Arrow Right - Skip to next frame</li>
                    <li>Arrow Left - Skip to previous frame</li>
                    <li>"," comma - Skip to next frames</li>
                    <li>"." dot - Skip to previous frames</li>
                </ul>
            </div>
            <div class="modal-footer">
                <button
                    disabled={videoIsBeingExported}
                    class="btn btn-primary"
                    data-bs-dismiss="modal">Close</button
                >
            </div>
        </div>
    </div>
</div>

{#if msgSuccess}
    <div class="w-100 position-relative">
        <div
            class="alert alert-success position-fixed top-0 right-0"
            role="alert"
        >
            {msgSuccess}
        </div>
    </div>
{/if}

{#if msgFailed}
    <div class="w-100 position-relative">
        <div
            class="alert alert-danger position-fixed top-0 right-0"
            role="alert"
        >
            {msgFailed}
        </div>
    </div>
{/if}

<div class="mt-5" />

<style>
    .selected {
        border: 1px solid red;
    }
    video {
        height: 530px;
        width: 100%;
        cursor: pointer;
    }

    svg {
        height: 24px;
        width: 24px;
    }

    span.spinner-border.spinner-border-sm {
        height: 24px;
        width: 24px;
    }

    .h-100px {
        height: 100px;
    }

    .h-150px {
        height: 150px;
    }

    .default-img {
        border-radius: 5px;
        width: 140px;
        object-fit: contain;
    }

    .right-0 {
        right: 0;
    }
</style>
