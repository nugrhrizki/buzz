import { ParentProps } from "solid-js";

import levitate from "@/assets/levitate.svg";

function AuthLayout(props: ParentProps) {
  return (
    <div class="relative h-[100dvh] flex items-center">
      <div class="absolute left-0 top-0 flex gap-x-2 p-4">
        <img src="/notified.svg" />
        <h1 class="font-bold text-lg transition motion-reduce:transition-none overflow-hidden">Notified</h1>
      </div>
      <div class="p-8 flex-1">{props.children}</div>
      <div class="relative hidden h-full p-10 text-white dark:border-r flex-1 lg:block">
        <div class="absolute inset-4 bg-sky-800 flex items-center justify-center flex-col rounded-3xl">
          <img src={levitate} class="lg:w-[15em] lg:h-[15em] xl:w-[20em] xl:h-[20em] 2xl:w-[25em] 2xl:h-[25em]" />
          <h2 class="lg:text-xl xl:text-2xl 2xl:text-4xl font-bold text-center">
            LEVITATE YOUR TECH
            <br />
            TO ANOTHER LEVEL
          </h2>
        </div>
      </div>
    </div>
  );
}

export default AuthLayout;
