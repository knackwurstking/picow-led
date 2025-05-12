type onChange = (
    ev: Event & { currentTarget: HTMLInputElement },
) => void | Promise<void>;

export function create(
    title: string,
    value: number,
    onChange: onChange,
): HTMLElement {
    const t = document.querySelector<HTMLTemplateElement>(
        `template[name="color-range-slider"]`,
    )!;

    const item = (
        t.content.cloneNode(true) as HTMLElement
    ).querySelector<HTMLElement>("*")!;

    return update(item, title, value, onChange);
}

export function update(
    item: HTMLElement,
    title: string,
    value: number,
    onChange: onChange,
): HTMLElement {
    const titleElement = item.querySelector<HTMLElement>(`.title`)!;
    titleElement.innerText = title;

    const rangeInput =
        item.querySelector<HTMLInputElement>(`input[type="range"]`)!;

    rangeInput.value = value.toString();
    rangeInput.oninput = () => {
        if (!numberInput) return;
        numberInput.value = rangeInput.value;
    };
    rangeInput.onchange = (ev) => {
        // @ts-expect-error
        onChange(ev);
    };

    const numberInput =
        item.querySelector<HTMLInputElement>(`input[type="number"]`)!;

    numberInput.value = value.toString();
    numberInput.onchange = (ev) => {
        rangeInput.value = numberInput.value;
        // @ts-expect-error
        onChange(ev);
    };

    return item;
}
