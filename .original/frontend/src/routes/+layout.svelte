<script lang="ts">
    import "./bootstrap-icons.css";

    import Toasts from "$lib/components/Toasts.svelte";

    let { children } = $props();
</script>

{@render children()}

<Toasts />
