export function urlQueryParam(name: string): string {
    if (!location.search) {
        return "";
    }

    for (const q of location.search.slice(1).split("&")) {
        const [k, v] = q.split("=", 2);
        if (k === name) {
            return decodeURIComponent(v);
        }
    }

    return "";
}

export function toHex(color: Color): string {
    const r = color.r.toString(16).padStart(2, "0");
    const g = color.g.toString(16).padStart(2, "0");
    const b = color.b.toString(16).padStart(2, "0");

    return `#${r}${g}${b}`;
}

export function toColor(value: string): Color {
    value = (value || "#FFFFFF").slice(1);
    const colorArray: number[] = [];

    for (let x = 0; x < value.length; x += 2) {
        colorArray.push(parseInt(value.slice(x, x + 2), 16));
    }

    const color: Color = {
        r: 0,
        g: 0,
        b: 0,
    };

    colorArray.slice(0, 3).forEach((c, i) => {
        switch (i) {
            case 0:
                color.r = c;
                break;
            case 1:
                color.g = c;
                break;
            case 2:
                color.b = c;
                break;
        }
    });

    return color;
}
