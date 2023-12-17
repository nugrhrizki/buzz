import { createForm, zodForm } from "@modular-forms/solid";
import { useNavigate } from "@solidjs/router";
import { TbLoader } from "solid-icons/tb";
import { Show } from "solid-js";

import { Button } from "@/components/ui/button";
import { Grid } from "@/components/ui/grid";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { type AuthForm, authFormSchema } from "../validations/auth";

export function LoginForm() {
  const navigate = useNavigate();
  const [authForm, { Form, Field }] = createForm<AuthForm>({
    validate: zodForm(authFormSchema),
  });

  function handleSubmit(auth: AuthForm) {
    return new Promise((resolve) =>
      setTimeout(() => {
        console.log(auth);
        navigate("/", { replace: true });
        resolve(true);
      }, 2000),
    );
  }

  return (
    <div class="grid gap-6">
      <Form onSubmit={handleSubmit}>
        <Grid class="gap-4">
          <Field name="username">
            {(field, props) => (
              <Grid class="gap-1">
                <Label for="username">Username</Label>
                <Input {...props} type="text" id="username" autocomplete="true" />
                <Show when={field.error}>
                  <p class="text-destructive text-xs">{field.error}</p>
                </Show>
              </Grid>
            )}
          </Field>
          <Field name="password">
            {(field, props) => (
              <Grid class="gap-1">
                <Label for="password">Password</Label>
                <Input {...props} type="password" id="password" />
                <Show when={field.error}>
                  <p class="text-destructive text-xs">{field.error}</p>
                </Show>
              </Grid>
            )}
          </Field>
          <Button type="submit" disabled={authForm.submitting} class="mt-8">
            <Show when={authForm.submitting}>
              <TbLoader class="mr-2 h-4 w-4 animate-spin" />
            </Show>
            Login
          </Button>
        </Grid>
      </Form>
    </div>
  );
}
