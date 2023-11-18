<script lang="ts">
    import { type WebResponse } from "../types/webResponse";

    let upload: SVGElement;
    let showSuccessDialog = false;
    let showFailedDialog = false;
    let uploadInProgress = false;

    function openUploadVideoDialog() {
        document.getElementById("uploadVideo").click();
        console.log("Opened upload video dialog");
    }

    function sendFileToServer(event: Event) {
        const inputElement = event.target as HTMLInputElement;
        if (inputElement.files && inputElement.files.length) {
            uploadInProgress = true;
            const file: File = inputElement.files[0];
            const formData = new FormData();
            formData.append("file", file);
            fetch("http://localhost:3000/upload-video", {
                method: "POST",
                body: formData,
            })
                .then((response) => response.json())
                .then((data: WebResponse) => {
                    console.log(data);

                    uploadInProgress = false;

                    if (data.Status == "success") {
                        showSuccessDialog = true;
                        document.location.reload();
                    }

                    if (data.Status == "failed") {
                        showSuccessDialog = true;
                    }

                    setTimeout(() => {
                        showSuccessDialog = false;
                        showFailedDialog = false;
                    }, 3000);
                });
        }
    }
</script>

<h1 class="text-white mt-5">SSCut</h1>
<p class="text-white">
    No videos yet uploaded. Upload one video and start cutting seconds.
</p>
<div class="d-flex flex-column justify-content-center" style="margin-top: 40%;">
    <button
        title="Upload video"
        class="btn btn-primary btn-lg"
        on:click={openUploadVideoDialog}
        disabled={uploadInProgress}
    >
        {#if uploadInProgress}
            <span
                class="spinner-border spinner-border-sm"
                role="status"
                aria-hidden="true"
            />
        {/if}

        <span class="me-2">Upload video</span>

        <svg
            bind:this={upload}
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            fill="currentColor"
            viewBox="0 0 16 16"
        >
            <path
                d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5z"
            />
            <path
                d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708l3-3z"
            />
        </svg>
    </button>
</div>

<input
    on:change={sendFileToServer}
    class="d-none"
    type="file"
    id="uploadVideo"
/>

{#if showSuccessDialog}
    <div class="w-100 position-relative">
        <div
            class="alert alert-success position-fixed top-0 right-0"
            role="alert"
        >
            Video file was uploaded!
        </div>
    </div>
{/if}

{#if showFailedDialog}
    <div class="w-100 position-relative">
        <div
            class="alert alert-danger position-fixed top-0 right-0"
            role="alert"
        >
            Failed to upload video!
        </div>
    </div>
{/if}

<style>
    .right-0 {
        right: 0;
    }
</style>
