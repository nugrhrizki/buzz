import { A } from "@solidjs/router";
import { TbShieldStar, TbUsers } from "solid-icons/tb";
import { For, Show } from "solid-js";

import { Menu } from "@/types";

import { Card, CardContent } from "@/components/ui/card";
import { Grid } from "@/components/ui/grid";

function config(): Menu[] {
  return [
    {
      name: "User",
      icon: <TbUsers class="w-5 h-5" />,
      href: "/user",
      show: true,
    },
    {
      name: "Role",
      icon: <TbShieldStar class="w-5 h-5" />,
      href: "/role",
      show: true,
    },
  ];
}

config.authorize = function () {
  return true;
};

function ConfigPage() {
  return (
    <div class="space-y-4 p-8 pt-6">
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight">Config</h2>
      </div>
      {/* grid 6 colums grid 4 colums 2 columns */}
      <Grid cols={1} colsMd={2} colsLg={4} class="w-full gap-4">
        <For each={config()}>
          {(item) => (
            <Show when={item.show && item.href}>
              <A href={"/config" + item.href!}>
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

export default ConfigPage;
export { config };
