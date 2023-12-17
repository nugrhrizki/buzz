import { z } from "zod";

export const authFormSchema = z.object({
  username: z.string().min(1, {
    message: "Username cannot be empty",
  }),
  password: z.string().min(1, {
    message: "Password cannot be empty",
  }),
});

export type AuthForm = z.infer<typeof authFormSchema>;
