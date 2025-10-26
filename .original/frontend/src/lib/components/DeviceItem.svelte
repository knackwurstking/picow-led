<script lang="ts">
    import SVGPower from "./SVGPower.svelte";

    import * as api from "$lib/api";

    let { addr, name, color, pins, active_color, power }: Device = $props();

    let powerButtonState = $derived<"off" | "on" | "processing">(
        power === 1 ? "on" : "off",
    );

    async function editButtonClick() {
        location.href = `./device/?addr=${addr}`;
    }

    async function powerButtonClick() {
        powerButtonState = "processing";

        await api.devices.addr.power.POST(
            addr,
            power === api.PowerStateON ? api.PowerStateOFF : api.PowerStateON,
        );
    }
</script>

<li class="ui-flex row gap justify-between align-center ui-padding">
    <span class="ui-padding ui-flex column justify-start">
        <h3>
            {#if name}
                {name}
            {:else}
                {addr}
            {/if}
        </h3>

        {#if name}
            <h6>
                {addr}
            </h6>
        {/if}
    </span>

    <span class="ui-flex-item ui-flex row gap" style="flex: 0">
        <span class="ui-flex-item" style="flex: 0">
            <button
                class="edit"
                data-ui-variant="ghost"
                data-ui-color="secondary"
                onclick={editButtonClick}
            >
                Edit
            </button>
        </span>

        <span class="ui-flex-item" style="flex: 0">
            <button
                class="power"
                data-ui-variant="ghost"
                data-ui-icon
                data-state={powerButtonState}
                onclick={powerButtonClick}
            >
                <div
                    class="background"
                    style="background-color: rgb({color.length >= 3
                        ? color.slice(0, 3).join(', ')
                        : '0, 0, 0'})"
                ></div>

                <SVGPower />
            </button>
        </span>
    </span>
</li>

<style>
    li {
        z-index: 0;
        border-radius: var(--ui-radius);

        width: 100%;

        backface-visibility: hidden;
        -webkit-backface-visibility: hidden;

        overflow: visible;
    }

    button.power {
        width: 3rem;
        height: 3rem;
        overflow: visible;
        color: rgb(0, 0, 0) !important;
    }

    button.power[data-state="on"] {
        color: rgb(0, 225, 0) !important;
    }

    button.power[data-state="processing"] {
        color: yellow !important;
    }

    button.power[data-state="off"] {
        color: rgb(0, 0, 0) !important;
    }

    button.power .background {
        position: absolute;
        z-index: -1;
        top: 0.25rem;
        right: 0.25rem;
        bottom: 0.25rem;
        left: 0.25rem;
        border-radius: var(--ui-radius);

        filter: blur(0.5rem);
        -webkit-backdrop-filter: blur(0.5rem);
    }
</style>
