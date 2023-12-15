import type { ComponentProps } from "solid-js";
import { splitProps } from "solid-js";

import { cn } from "@/lib/utils";

export function MainNav(props: ComponentProps<"nav">) {
  const [, rest] = splitProps(props, ["class"]);
  return (
    <nav class={cn("flex items-center space-x-4 lg:space-x-6", props.class)} {...rest}>
      <a href="/examples/dashboard" class="hover:text-primary text-sm font-medium transition-colors">
        Overview
      </a>
      <a
        href="/examples/dashboard"
        class="text-muted-foreground hover:text-primary text-sm font-medium transition-colors">
        Customers
      </a>
      <a
        href="/examples/dashboard"
        class="text-muted-foreground hover:text-primary text-sm font-medium transition-colors">
        Products
      </a>
      <a
        href="/examples/dashboard"
        class="text-muted-foreground hover:text-primary text-sm font-medium transition-colors">
        Settings
      </a>
    </nav>
  );
}
