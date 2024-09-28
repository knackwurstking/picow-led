import { moreVertical as svgOptions } from "ui/svg/smoothie-line-icons";
import { html, UIIconButton } from "ui";

export default class PicowOptionsButton extends UIIconButton {
    constructor() {
        super();
        this.picow = {};
        this.#render();
    }

    #render() {
        this.ui.ghost = true;

        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    height: 100%;
                }
            </style>
        `;

        this.innerHTML = svgOptions;

        this.onclick = async (ev) => {
            ev.stopPropagation();

            // TODO: Open the device setup dialog
        };
    }
}

console.debug(`Register the "picow-options-button"`);
customElements.define("picow-options-button", PicowOptionsButton);
