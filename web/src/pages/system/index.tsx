import { A } from "@solidjs/router";
import { TbFlag, TbHistory, TbSettings } from "solid-icons/tb";
import { For, Show } from "solid-js";

import { Menu } from "@/types";

import { Card, CardContent } from "@/components/ui/card";
import { Grid } from "@/components/ui/grid";

function system(): Menu[] {
  return [
    {
      name: "Flag",
      icon: <TbFlag class="w-5 h-5" />,
      href: "/flag",
      show: true,
    },
    {
      name: "Log",
      icon: <TbHistory class="w-5 h-5" />,
      href: "/log",
      show: true,
    },
    {
      name: "Setting",
      icon: <TbSettings class="w-5 h-5" />,
      href: "/setting",
      show: true,
    },
  ];
}

system.authorize = function () {
  return false;
};

function SystemPage() {
  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">System</h2>
      </div>
      {/* grid 6 colums grid 4 colums 2 columns */}
      <Grid cols={1} colsMd={2} colsLg={4} class="w-full gap-4">
        <For each={system()}>
          {(item) => (
            <Show when={item.show && item.href}>
              <A href={"/system" + item.href!}>
                <Card>
                  <CardContent class="pt-6 flex items-center gap-3 text-xl">
                    {item.icon}
                    <span>{item.name}</span>
                  </CardContent>
                </Card>
              </A>
            </Show>
          )}
        </For>
      </Grid>
    </div>
  );
}

export default SystemPage;
export { system };
