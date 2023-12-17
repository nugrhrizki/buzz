import { As, useColorMode } from "@kobalte/core";
import { TbDeviceLaptop, TbMoon, TbSun } from "solid-icons/tb";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function ModeToggle() {
  const { setColorMode } = useColorMode();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <As component={Button} variant="ghost" size="sm">
          <TbSun class="w-5 h-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <TbMoon class="w-5 h-5 absolute rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          <span class="sr-only">Toggle theme</span>
        </As>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="z-[99]">
        <DropdownMenuItem onSelect={() => setColorMode("light")}>
          <TbSun class="mr-2 h-4 w-4" />
          <span>Light</span>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setColorMode("dark")}>
          <TbMoon class="mr-2 h-4 w-4" />
          <span>Dark</span>
        </DropdownMenuItem>
        <DropdownMenuItem onSelect={() => setColorMode("system")}>
          <TbDeviceLaptop class="mr-2 h-4 w-4" />
          <span>System</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
