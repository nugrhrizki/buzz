import { As } from "@kobalte/core";
import { useNavigate } from "@solidjs/router";
import { createEffect } from "solid-js";

import { generateAvatarUrl } from "@/lib/utils";

import { useUser } from "@/services/auth";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

export function UserNav() {
  const user = useUser();
  const navigate = useNavigate();

  createEffect(() => {
    if (user.status === "error") {
      if ((user.error as unknown as Response | undefined)?.status === 401) {
        navigate("/auth", { replace: true });
      }
    }
  });

  return (
    <DropdownMenu placement="bottom-end">
      <DropdownMenuTrigger asChild>
        <As component={Button} variant="ghost" class="relative h-8 w-8 rounded-full">
          <Avatar class="h-8 w-8">
            <AvatarImage
              src={generateAvatarUrl({
                name: user.data?.data.username || "johndoe",
                style: "big-smile",
                backgroundColors: ["#d6e6ff", "#d7f9f8", "#ffffea", "#fff0d4", "#fbe0e0", "#e5d4ef"],
              })}
              alt={user.data?.data.username || "johndoe"}
            />
            <AvatarFallback>{user.data?.data.username.charAt(0).toUpperCase() || "J"}</AvatarFallback>
          </Avatar>
        </As>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="w-56 z-[100]">
        <DropdownMenuLabel class="font-normal">
          <div class="flex flex-col space-y-1">
            <p class="text-sm font-medium leading-none">{user.data?.data.name}</p>
            <p class="text-muted-foreground text-xs leading-none">{user.data?.data.username}</p>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onSelect={() => {
            navigate("/auth", { replace: true });
          }}>
          Log out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
