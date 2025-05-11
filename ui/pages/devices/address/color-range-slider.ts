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
    );
    if (!t) {
        throw new Error(
            `Nope, template with name "color-range-slider" is null`,
        );
    }

    const item = (
        t.content.cloneNode(true) as HTMLElement
    ).querySelector<HTMLElement>("*");
    if (!item) throw new Error(`template is empty`);

    return update(item, title, value, onChange);
}

export function update(
    item: HTMLElement,
    title: string,
    value: number,
    onChange: onChange,
): HTMLElement {
    const titleElement = item.querySelector<HTMLElement>(`.title`);
    if (!titleElement) throw new Error(`title element is null`);
    titleElement.innerText = title;

    const rangeInput =
        item.querySelector<HTMLInputElement>(`input[type="range"]`);
    if (!rangeInput) throw new Error(`range input element is null`);

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
        item.querySelector<HTMLInputElement>(`input[type="number"]`);
    if (!numberInput) throw new Error(`number input element is null`);

    numberInput.value = value.toString();
    numberInput.onchange = (ev) => {
        rangeInput.value = numberInput.value;
        // @ts-expect-error
        onChange(ev);
    };

    return item;
}
