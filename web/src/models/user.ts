import { z } from "zod";

export const userSchema = z.object({
  id: z.number().optional().nullable().nullish(),
  name: z.string(),
  username: z.string(),
  password: z.string(),
  confirmed: z.boolean().optional().nullable().nullish(),
  whatsapp: z.string().optional().nullable().nullish(),
  email: z.string().optional().nullable().nullish(),
  role_id: z.number().optional().nullable().nullish(),
  created_at: z.date().optional().nullable().nullish(),
  updated_at: z.date().optional().nullable().nullish(),
  deleted_at: z.date().optional().nullable().nullish(),
});

export const createUserSchema = userSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
  deleted_at: true,
});

export type User = z.infer<typeof userSchema>;
export type CreateUser = z.infer<typeof createUserSchema>;
