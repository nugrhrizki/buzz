/* @refresh reload */
import { ColorModeProvider, ColorModeScript } from "@kobalte/core";
import { render } from "solid-js/web";

import "./root.css";
import { Routes } from "./routes";

render(
  () => (
    <>
      <ColorModeScript storageType="localStorage" />
      <ColorModeProvider>
        <Routes />
      </ColorModeProvider>
    </>
  ),
  document.getElementById("root")!,
);
