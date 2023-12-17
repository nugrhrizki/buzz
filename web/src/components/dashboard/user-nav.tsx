import { As } from "@kobalte/core";
import { useNavigate } from "@solidjs/router";

import { generateAvatarUrl } from "@/lib/utils";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function UserNav() {
  const navigate = useNavigate();

  return (
    <DropdownMenu placement="bottom-end">
      <DropdownMenuTrigger asChild>
        <As component={Button} variant="ghost" class="relative h-8 w-8 rounded-full">
          <Avatar class="h-8 w-8">
            <AvatarImage
              src={generateAvatarUrl({
                name: "nugrhrizki",
                style: "big-smile",
                backgroundColors: ["#d6e6ff", "#d7f9f8", "#ffffea", "#fff0d4", "#fbe0e0", "#e5d4ef"],
              })}
              alt="@mangiki"
            />
            <AvatarFallback>MI</AvatarFallback>
          </Avatar>
        </As>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="w-56 z-[100]">
        <DropdownMenuLabel class="font-normal">
          <div class="flex flex-col space-y-1">
            <p class="text-sm font-medium leading-none">nugrhrizki</p>
            <p class="text-muted-foreground text-xs leading-none">me@mangiki.com</p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onSelect={() => {
            navigate("/auth", { replace: true });
          }}>
          Log out
          <DropdownMenuShortcut>⇧⌘Q</DropdownMenuShortcut>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
