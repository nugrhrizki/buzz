import type { ParentProps } from "solid-js";

function Header(props: ParentProps) {
  return (
    <div class="border-b border-border sticky top-0 z-[98] bg-background/80 backdrop-blur-lg">{props.children}</div>
  );
}

export default Header;
