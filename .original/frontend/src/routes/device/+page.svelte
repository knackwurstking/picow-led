<script lang="ts">
    import { onMount } from "svelte";
    import { WS } from "ui";

    import { api, utils } from "$lib";
    api.setPathPrefix("..");

    import ColorRangeSlider from "$lib/components/ColorRangeSlider.svelte";
    import OnlineIndicator from "$lib/components/OnlineIndicator.svelte";
    import ColorStorageItem from "$lib/components/ColorStorageItem.svelte";
    import { PowerStateOFF, PowerStateON } from "$lib/api";

    const queryAddr = utils.urlQueryParam("addr");

    let device = $state<Device | undefined>(undefined);
    let onlineIndicator_DataState = $state<"offline" | "online">(
        window.ws?.isOpen() ? "online" : "offline",
    );

    let colors = $state<Colors | undefined>(undefined);
    let activeColorIndex = $state<number>(-1);

    async function powerOFF() {
        if (!device?.addr) {
            return;
        }

        await api.devices.addr.power.POST(device.addr, PowerStateOFF);
    }

    async function powerON() {
        if (!device?.addr) {
            return;
        }

        await api.devices.addr.power.POST(device.addr, PowerStateON);
    }

    async function addNewColorStorageItem(hexColor: string) {
        if (!device?.addr) {
            return;
        }

        const color = utils.toColor(hexColor);

        if (!colors) {
            colors = [color];
        } else {
            colors.push(color);
        }

        await api.colors.POST([color]);
    }

    onMount(async () => {
        window.ws = new WS<WSMessageData>(`../ws`, true);

        await window.ws.connect();

        window.ws.events.addListener("open", async () => {
            onlineIndicator_DataState = "online";

            device = await api.devices.addr.GET(queryAddr);
            updateActiveColorIndex();
        });

        window.ws.events.addListener("close", () => {
            onlineIndicator_DataState = "offline";
        });

        // Handle WebSocket message events ("devices", "device")
        window.ws.events.addListener("message", async (data) => {
            switch (data.type) {
                case "colors":
                    colors = data.data;
                    updateActiveColorIndex();
                    break;

                case "device":
                    {
                        if (device?.addr !== data.data.addr) {
                            return;
                        }

                        console.debug(`ws "device" event:`, data.data);
                        device = data.data;
                    }
                    break;
            }
        });

        colors = await api.colors.GET();
        updateActiveColorIndex();
    });

    function updateActiveColorIndex() {
        if (!colors || !device?.active_color) return;

        const match =
            colors.find((c) => {
                if (
                    c.r === device!.active_color[0] &&
                    c.g === device!.active_color[1] &&
                    c.b === device!.active_color[2]
                ) {
                    return true;
                }

                return false;
            }) || -1;

        if (match !== -1) {
            activeColorIndex = colors.indexOf(match);
        }
    }
</script>

<svelte:head>
    <title>PicoW LED | {queryAddr}</title>
</svelte:head>

<main class="ui-container">
    <div class="ui-flex column gap">
        <span class="power ui-flex gap justify-between align-center">
            <span>Power</span>

            <span>
                <button
                    class="off"
                    data-ui-color="destructive"
                    onclick={powerOFF}
                >
                    OFF
                </button>

                <button class="on" onclick={powerON}>ON</button>
            </span>
        </span>

        <h2>Color Storage</h2>

        {#if device && device.pins.length > 3}
            <div class="range-sliders ui-flex column gap nowrap">
                {#each device.pins.slice(3) as pin, index}
                    <ColorRangeSlider
                        {pin}
                        bind:value={device.color[index + 3]}
                        onchange={async () => {
                            if (!device || !device.addr || !device.color) {
                                return;
                            }

                            await api.devices.addr.color.POST(
                                device.addr,
                                device.color,
                            );
                        }}
                    />
                {/each}
            </div>
        {/if}

        <br />

        {#if colors && colors.length > 0}
            <div class="color-storage-container ui-flex row gap wrap">
                {#each colors as color, index}
                    <ColorStorageItem
                        {color}
                        active={activeColorIndex === index}
                        onclick={async () => {
                            if (activeColorIndex === index) {
                                return;
                            }

                            activeColorIndex = index;

                            if (device?.addr) {
                                const newColor = [
                                    color.r,
                                    color.g,
                                    color.b,
                                    ...device.color.slice(3),
                                ];

                                await api.devices.addr.color.POST(
                                    device.addr,
                                    newColor,
                                );
                            }
                        }}
                        onchange={async (newColor) => {
                            console.debug("Color changed:", color.id, newColor);

                            colors![index] = { ...color, ...newColor };

                            if (color.id !== undefined) {
                                await api.colors.id.POST(color.id, newColor);
                            }
                        }}
                        ondelete={async () => {
                            if (!colors) return;

                            api.colors.PUT([
                                ...colors.filter((c) => c.id !== color.id),
                            ]);
                        }}
                    />
                {/each}
            </div>
        {/if}

        <button
            class="new-color ui-flex-item"
            onclick={() => {
                addNewColorStorageItem(
                    utils.toHex({
                        r: 255,
                        g: 255,
                        b: 255,
                    }),
                );
            }}
        >
            New Color
        </button>
    </div>
</main>

<div class="ui-app-bar" data-ui-position="top">
    <span class="ui-app-bar-left">
        <OnlineIndicator state={onlineIndicator_DataState} />
    </span>

    <span class="ui-app-bar-center">
        <h4>{device?.name || queryAddr}</h4>
    </span>

    <span class="ui-app-bar-right"></span>
</div>

<style>
    main {
        padding-top: calc(var(--ui-app-bar-height) + var(--ui-spacing));
    }
</style>
