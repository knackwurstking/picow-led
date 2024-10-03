interface AppBar {
    element: import("ui").UIAppBar;
    events: import("ui").Events<AppBar_Events>;
    items: {
        menu: import("ui").UIAppBarItem<import("ui").UIIconButton>;
        title: import("ui").UIAppBarItem<HTMLElement>;
        add: import("ui").UIAppBarItem<import("ui").UIIconButton>;
    };
    buttons: {
        menu: import("ui").UIIconButton;
        add: import("ui").UIIconButton;
    };
    title: string;
}

interface AppBar_Events {
    menu: MouseEvent & { currentTarget: import("ui").UIIconButton };
    add: MouseEvent & { currentTarget: import("ui").UIIconButton };
}
