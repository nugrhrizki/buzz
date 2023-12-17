import { sidebarPin, toggleSidebarPin } from "@/store/global";
import { useColorMode } from "@kobalte/core";
import { useNavigate } from "@solidjs/router";
import { CgMenu, CgMenuLeftAlt } from "solid-icons/cg";
import { TbSearch } from "solid-icons/tb";
import { ParentProps, Show } from "solid-js";
// @ts-ignore
import { NinjaKeys, createNinjaKeys } from "solid-ninja-keys";

import { actions } from "@/actions";

import { Button } from "@/components/ui/button";

import { ModeToggle } from "@/components/dark-theme";
import { UserNav } from "@/components/dashboard/user-nav";
import Header from "@/components/header";
import Navigation, { menu } from "@/components/navigation";
import Sidebar from "@/components/sidebar";

function AdminLayout(props: ParentProps) {
  const { open } = createNinjaKeys();
  const color = useColorMode();
  const navigate = useNavigate();

  return (
    <>
      <Sidebar>
        <Navigation />
      </Sidebar>
      <main
        class="transition-[margin] motion-reduce:transition-none"
        classList={{ "ml-[250px]": sidebarPin(), "ml-[70px]": !sidebarPin() }}>
        <Header>
          <div class="flex h-16 items-center gap-x-4 px-4">
            <Button variant="ghost" size="icon" onClick={toggleSidebarPin}>
              <Show when={sidebarPin()} fallback={<CgMenu />}>
                <CgMenuLeftAlt class="w-5 h-5" />
              </Show>
            </Button>
            <Button variant="ghost" size="icon" onClick={() => open()}>
              <TbSearch class="w-5 h-5" />
            </Button>
            <div class="ml-auto flex items-center space-x-4">
              <ModeToggle />
              <UserNav />
            </div>
          </div>
        </Header>
        {props.children}
      </main>
      <NinjaKeys
        hotkeys={actions({
          color,
          navigate,
          menu,
        })}
        isDark={color.colorMode() === "dark"}
      />
    </>
  );
}

export default AdminLayout;
