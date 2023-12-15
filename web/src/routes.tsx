import { Route, HashRouter } from "@solidjs/router";
import { lazy } from "solid-js";

import AdminLayout from "@/layout/admin";
import AuthLayout from "@/layout/auth";
import BaseLayout from "@/layout/base";

const AuthPage = lazy(() => import("@/pages/auth"));
const DashboardPage = lazy(() => import("@/pages/dashboard"));
const NotFoundPage = lazy(() => import("@/pages/404"));

export function Routes() {
  return (
    <HashRouter root={BaseLayout}>
      <Route path="/" component={AdminLayout}>
        <Route path="/" component={DashboardPage} />
      </Route>
      <Route path="/auth" component={AuthLayout}>
        <Route path="/" component={AuthPage} />
      </Route>
      <Route path="/*" component={NotFoundPage} />
    </HashRouter>
  );
}
