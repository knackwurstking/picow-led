<script lang="ts">
    // NOTE: https://svelte.dev/playground/0091c8b604b74ed88bb7b6d174504f50?version=3.35.0

    import { toasts, removeToast } from "../stores/toasts";

    import Toast from "./Toast.svelte";
</script>

<section
    class="toasts ui-flex column gap justify-end align-end ui-auto-scroll ui-hide-scrollbar"
>
    {#each $toasts as toast}
        <Toast
            type={toast.type}
            timeout={toast.timeout}
            ondismiss={async () => {
                removeToast(toast.id);
            }}
        >
            {toast.message}
        </Toast>
    {/each}
</section>

<style>
    .toasts {
        z-index: 999;
        position: absolute;
        bottom: var(--ui-spacing);
        right: var(--ui-spacing);
        width: fit-content;
        max-width: calc(100% - var(--ui-spacing) * 2);
        height: fit-content;
        max-height: calc(100% - var(--ui-spacing) * 2);
    }
</style>
