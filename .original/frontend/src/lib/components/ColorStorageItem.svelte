<script lang="ts">
    import { toColor, toHex } from "$lib/utils";

    let {
        color = $bindable(),
        active = false,
        onclick,
        onchange,
        ondelete,
    }: {
        color: Color;
        active: boolean;
        onclick?: (() => void | Promise<void>) | null;
        onchange?: ((color: Color) => void | Promise<void>) | null;
        ondelete?: (() => void | Promise<void>) | null;
    } = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class="color-storage-item ui-border ui-padding"
    style:color={`rgb(${color.r}, ${color.g}, ${color.b})`}
    class:active
    {onclick}
>
    <div class="background"></div>

    <input
        type="color"
        value={toHex(color)}
        onchange={async (ev) => {
            color = toColor(ev.currentTarget.value);
            if (onchange) {
                onchange(color);
            }
        }}
    />

    <!-- svelte-ignore a11y_consider_explicit_label -->
    <button
        class="delete"
        data-ui-variant="ghost"
        data-ui-icon
        onclick={(ev) => {
            ev.stopPropagation();

            if (ondelete) {
                ondelete();
            }
        }}
    >
        <i class="bi bi-x"></i>
    </button>
</div>

<style>
    .color-storage-item {
        width: 7rem;
        height: 3.5rem;
        outline: var(--ui-border-width) solid transparent;
        outline-offset: var(--ui-border-width);
        border-radius: var(--ui-radius);
        cursor: pointer;
    }

    .color-storage-item .background {
        background-color: currentColor;
        color: inherit;
        width: 100%;
        height: 100%;
        border-radius: inherit;
    }

    .color-storage-item input[type="color"] {
        display: none;
        opacity: 0;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }

    .color-storage-item button.delete {
        --wght: 850;
        position: absolute;
        top: 0;
        right: 0;
        transform: translate(50%, -50%);
        border-radius: 50%;
        padding: 0 !important;
    }

    .color-storage-item.active {
        outline: var(--ui-border-width) solid var(--ui-primary);
    }

    .color-storage-item.active input[type="color"] {
        display: block;
        position: absolute;
    }
</style>
