/* @refresh reload */
import { render } from "solid-js/web";

import "./root.css";
import { Routes } from "./routes";

render(() => <Routes />, document.getElementById("root")!);
