import { HashRouter, Route } from "@solidjs/router";
import { lazy } from "solid-js";

import AdminLayout from "@/layout/admin";
import AuthLayout from "@/layout/auth";
import BaseLayout from "@/layout/base";

const LoginPage = lazy(() => import("@/pages/login"));

const DashboardPage = lazy(() => import("@/pages/dashboard"));
const ClientPage = lazy(() => import("@/pages/client"));

const ConfigPage = lazy(() => import("@/pages/config"));
const UserPage = lazy(() => import("@/pages/config/user"));
const RolePage = lazy(() => import("@/pages/config/role"));

const SystemPage = lazy(() => import("@/pages/system"));
const FlagPage = lazy(() => import("@/pages/system/flag"));
const LogPage = lazy(() => import("@/pages/system/log"));
const SettingPage = lazy(() => import("@/pages/system/setting"));

const NotFoundPage = lazy(() => import("@/pages/404"));

export function Routes() {
  return (
    <HashRouter root={BaseLayout}>
      <Route path="/" component={AdminLayout}>
        <Route path="/" component={DashboardPage} />
        <Route path="/client" component={ClientPage} />

        <Route path="/config" component={ConfigPage} />
        <Route path="/config/user" component={UserPage} />
        <Route path="/config/role" component={RolePage} />

        <Route path="/system" component={SystemPage} />
        <Route path="/system/flag" component={FlagPage} />
        <Route path="/system/log" component={LogPage} />
        <Route path="/system/setting" component={SettingPage} />
      </Route>
      <Route path="/auth" component={AuthLayout}>
        <Route path="/" component={LoginPage} />
      </Route>
      <Route path="/*" component={NotFoundPage} />
    </HashRouter>
  );
}
