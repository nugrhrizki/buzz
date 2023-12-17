import { sidebarPin } from "@/store/global";
import { A } from "@solidjs/router";
import { TbBrandWhatsapp, TbDashboard } from "solid-icons/tb";
import { For, Match, Show, Switch } from "solid-js";

import { Menu } from "@/types";

import { config } from "@/pages/config";
import { system } from "@/pages/system";

import NavLink from "@/components/navlink";

export function menu(): Menu[] {
  return [
    {
      name: "Menu",
      isGroup: true,
      show: true,
      children: [
        {
          name: "Dashboard",
          icon: <TbDashboard class="w-5 h-5" />,
          href: "/",
          show: true,
        },
        {
          name: "Sender",
          icon: <TbBrandWhatsapp class="w-5 h-5" />,
          href: "/sender",
          show: true,
        },
      ],
    },
    {
      name: "Config",
      isGroup: true,
      show: config.authorize(),
      href: "/config",
      children: config(),
    },
    {
      name: "System",
      isGroup: true,
      show: system.authorize(),
      href: "/system",
      children: system(),
    },
  ];
}

function Navigation() {
  return (
    <For each={menu()}>
      {(item) => (
        <Switch>
          <Match when={item.isGroup && item.show}>
            <div
              class="flex flex-col gap-y-2"
              classList={{
                "my-8": sidebarPin(),
                "group-hover:my-8": !sidebarPin(),
              }}>
              <Show when={item.href !== undefined}>
                <A
                  href={item.href!}
                  class="ml-4 font-bold text-medium group-hover:h-[initial] overflow-hidden transition-[height] motion-reduce:transition-none hover:underline"
                  classList={{
                    "h-0": !sidebarPin(),
                  }}>
                  {item.name}
                </A>
              </Show>
              <Show when={item.href === undefined}>
                <span
                  class="ml-4 font-bold text-medium group-hover:h-[initial] overflow-hidden transition-[height] motion-reduce:transition-none"
                  classList={{
                    "h-0": !sidebarPin(),
                  }}>
                  {item.name}
                </span>
              </Show>
              <For each={item.children}>
                {(child) => (
                  <Show when={child.show && child.href !== undefined}>
                    <NavLink href={item.href !== undefined ? item.href! + child.href! : child.href!}>
                      {child.icon}
                      <span class="font-medium overflow-hidden">{child.name}</span>
                    </NavLink>
                  </Show>
                )}
              </For>
            </div>
          </Match>
          <Match when={!item.isGroup && item.show && item.href !== undefined}>
            <NavLink href={item.href!}>
              {item.icon}
              <span class="font-medium overflow-hidden">{item.name}</span>
            </NavLink>
          </Match>
        </Switch>
      )}
    </For>
  );
}

export default Navigation;
