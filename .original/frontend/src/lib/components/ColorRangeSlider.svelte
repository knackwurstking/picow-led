<script lang="ts">
    let {
        value = $bindable(),
        pin,
        onchange,
    }: {
        value: number;
        pin: number;
        onchange: () => void | Promise<void>;
    } = $props();
</script>

<span class="color-range-slider ui-flex column justify-start">
    <span class="title">{pin}</span>

    <span class="container ui-flex row gap nowrap">
        <input type="range" min="0" max="255" step="1" bind:value {onchange} />
        <input type="number" min="0" max="255" step="1" bind:value {onchange} />
    </span>
</span>

<style>
    .color-range-slider {
        height: fit-content;
        width: 100%;
    }

    .color-range-slider > .title {
        --slnt: -15;
        width: fit-content;
        font-size: 0.9rem;
    }

    .color-range-slider > .container {
        width: 100%;
    }

    .color-range-slider > .container > input[type="number"] {
        height: 100%;
        width: 6ch;
        text-align: center;
    }
</style>
