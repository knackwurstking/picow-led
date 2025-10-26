import { writable } from "svelte/store";

export type ToastDataTypes = "info" | "warning" | "error";

export interface ToastData {
    id: number;
    type: ToastDataTypes;
    message: string;
    timeout: number;
}

export const toasts = writable<ToastData[]>([]);

export function addToast(
    type: ToastDataTypes,
    message: string,
    timeout: number = 0,
): void {
    const id = new Date().getTime();

    toasts.update((storeData) => [
        ...storeData,
        { id, type, message, timeout },
    ]);

    if (timeout > 0) {
        setTimeout(() => removeToast(id), timeout);
    }
}

export function removeToast(id: number): void {
    toasts.update((storeData) => storeData.filter((t) => t.id !== id));
}
