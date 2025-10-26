<script lang="ts">
    import { onMount } from "svelte";
    import { WS } from "ui";

    import { api } from "$lib";
    api.setPathPrefix(".");

    import DeviceItem from "$lib/components/DeviceItem.svelte";
    import OnlineIndicator from "$lib/components/OnlineIndicator.svelte";

    let devices = $state<Devices>([]);
    let onlineIndicator_DataState = $state<"offline" | "online">(
        window.ws?.isOpen() || false ? "online" : "offline",
    );

    onMount(async () => {
        window.ws = new WS<WSMessageData>(`./ws`, true);
        await window.ws.connect();

        window.ws.events.addListener("open", async () => {
            onlineIndicator_DataState = "online";

            // Fetch devices from /api and render items
            devices = await api.devices.GET();
        });

        window.ws.events.addListener("close", () => {
            onlineIndicator_DataState = "offline";
        });

        // Handle WebSocket message events ("devices", "device")
        window.ws.events.addListener("message", async (data) => {
            switch (data.type) {
                case "devices":
                    {
                        console.debug(`ws "devices" event:`, data.data);

                        devices = data.data;
                    }
                    break;

                case "device":
                    {
                        console.debug(`ws "device" event:`, data.data);

                        for (let x = 0; x < devices.length; x++) {
                            const device = devices[x];
                            if (device.addr === data.data.addr) {
                                devices[x] = data.data;
                            }
                        }
                    }
                    break;
            }
        });
    });
</script>

<svelte:head>
    <title>PicoW LED | Devices</title>
</svelte:head>

<main class="ui-container">
    <ul>
        {#each devices as device}
            <DeviceItem {...device} />
        {/each}
    </ul>
</main>

<div class="ui-app-bar" data-ui-position="top">
    <span class="ui-app-bar-left">
        <OnlineIndicator state={onlineIndicator_DataState} />
    </span>

    <span class="ui-app-bar-center">
        <h4>Devices</h4>
    </span>

    <span class="ui-app-bar-right"></span>
</div>

<style>
    main {
        padding-top: calc(var(--ui-app-bar-height) + var(--ui-spacing));
    }
</style>
