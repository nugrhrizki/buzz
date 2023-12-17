import { sidebarPin } from "@/store/global";
import { A } from "@solidjs/router";
import { type ParentProps } from "solid-js";

function Sidebar(props: ParentProps) {
  return (
    <div
      class="fixed h-[100dvh] group z-[99] border-border border-r hover:w-[250px] hover:p-4 bg-background text-foreground flex flex-col justify-between transition-[width] motion-reduce:transition-none"
      classList={{
        "w-[250px] p-4": sidebarPin(),
        "w-[70px] px-2 py-4": !sidebarPin(),
      }}>
      <div class="space-y-8">
        <A href="/" class="flex gap-x-2 px-4">
          <img src="/notified.svg" />
          <h1 class="font-bold text-lg transition motion-reduce:transition-none overflow-hidden">Notified</h1>
        </A>
        <div class="transition-[margin] motion-reduce:transition-none">{props.children}</div>
      </div>
      <span
        class="text-xs"
        classList={{
          "opacity-0 group-hover:opacity-100": !sidebarPin(),
          "opacity-100": sidebarPin(),
        }}>
        v0.1.0
      </span>
    </div>
  );
}

export default Sidebar;
