import * as globals from "./globals";

const colorStringSeparator = ",";

export const colorCache = {
    add(color: number[]): void {
        globals.store.update("color", (data) => {
            const newColorString = color.join(colorStringSeparator);
            if (data.cache.findIndex((c) => c.join(colorStringSeparator) === newColorString) >= 0) {
                return data;
            }

            data.cache.unshift(color);
            if (data.cache.length > 20) {
                data.cache = data.cache.slice(0, 20);
            }

            return data;
        });
    },

    getAll() {
        return globals.store.get("color")?.cache || [];
    },

    remove(color: number[]) {
        globals.store.update("color", (data) => {
            const colorString = color.join(colorStringSeparator);
            data.cache = data.cache.filter((c) => c.join(colorStringSeparator) !== colorString);
            return data;
        });
    },
};
