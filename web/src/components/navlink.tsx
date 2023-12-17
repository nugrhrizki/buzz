import { A, AnchorProps } from "@solidjs/router";

import { cn } from "@/lib/utils";

import { buttonVariants } from "@/components/ui/button";

function NavLink(props: AnchorProps) {
  return (
    <A
      activeClass="bg-secondary"
      class={cn(
        buttonVariants({
          variant: "ghost",
        }),
        "flex items-center justify-start space-x-2 hover:bg-secondary/80",
      )}
      {...props}
      end
    />
  );
}

export default NavLink;
