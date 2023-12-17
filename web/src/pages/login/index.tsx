import { LoginForm } from "./forms/login";

function LoginPage() {
  return (
    <div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
      <div class="flex flex-col space-y-2 mb-8">
        <h1 class="text-2xl font-semibold tracking-tight">Sign In</h1>
        <p class="text-muted-foreground text-sm">Enter your username and password to continue</p>
      </div>
      <LoginForm />
    </div>
  );
}

export default LoginPage;
