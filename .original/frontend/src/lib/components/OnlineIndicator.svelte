<script lang="ts">
    const {
        state
    }: {
        state: "offline" | "online"
    } = $props()
</script>

<span class="online-indicator" data-state={state}></span>

<style>
    .online-indicator {
        --mono: 1;
        --wght: 650;
        font-size: 0.75rem;
    }

    .online-indicator[data-state="online"] {
        color: green;
    }

    .online-indicator[data-state="online"]::after {
        content: "Online";
        color: green;
    }

    .online-indicator[data-state="offline"] {
        color: red;
    }

    .online-indicator[data-state="offline"]::after {
        content: "Offline";
        color: red;
    }
</style>
