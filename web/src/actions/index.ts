import { toggleSidebarPin } from "@/store/global";
import { ColorModeContextType } from "@kobalte/core";
import { Navigator } from "@solidjs/router";

import { Menu } from "@/types";

export type ActionCtx = {
  color: ColorModeContextType;
  navigate: Navigator;
  menu: () => Menu[];
};

function menuAction(ctx: ActionCtx) {
  const menus: Menu[] = [];

  for (const menu of ctx.menu()) {
    if (menu.href && menu.show) {
      menus.push(menu);
    }
    if (menu.children && menu.show) {
      for (const child of menu.children) {
        if (child.href && child.show) {
          if (menu.href) {
            child.href = menu.href + child.href;
          }
          menus.push(child);
        }
      }
    }
  }

  const actions = [];
  const children = [];

  for (const menu of menus) {
    const id = "Menu " + (menu.href?.replace(/\//g, " ").trim() || "Dashboard");
    children.push(id);
    actions.push({
      id,
      parent: "Menu",
      title: "Menu " + menu.name,
      handler: () => {
        if (menu.href) {
          ctx.navigate(menu.href);
        }
      },
    });
  }

  return [
    {
      id: "Menu",
      title: "Menu",
      mdIcon: "menu",
      children,
    },
    ...actions,
  ];
}

const toggleSidebarAction = {
  id: "Toggle Sidebar",
  title: "Toggle Sidebar",
  mdIcon: "toggle_on",
  hotkey: "ctrl+b",
  handler: toggleSidebarPin,
};

function switchThemeAction(ctx: ActionCtx) {
  return [
    {
      id: "Theme",
      title: "Change theme...",
      mdIcon: "desktop_windows",
      children: ["Light Theme", "Dark Theme", "System Theme"],
    },
    {
      id: "Light Theme",
      title: "Change theme to Light",
      mdIcon: "light_mode",
      parent: "Theme",
      handler: () => {
        ctx.color.setColorMode("light");
      },
    },
    {
      id: "Dark Theme",
      title: "Change theme to Dark",
      mdIcon: "dark_mode",
      parent: "Theme",
      handler: () => {
        ctx.color.setColorMode("dark");
      },
    },
    {
      id: "System Theme",
      title: "Change theme to System",
      mdIcon: "brightness_auto",
      parent: "Theme",
      handler: () => {
        ctx.color.setColorMode("system");
      },
    },
  ];
}

function logoutAction(ctx: ActionCtx) {
  return {
    id: "Logout",
    title: "Logout",
    mdIcon: "logout",
    handler: () => {
      ctx.navigate("/auth");
    },
  };
}

export function actions(ctx: ActionCtx) {
  return [...menuAction(ctx), toggleSidebarAction, ...switchThemeAction(ctx), logoutAction(ctx)];
}
