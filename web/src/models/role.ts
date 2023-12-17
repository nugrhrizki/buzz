import { z } from "zod";

export const roleSchema = z.object({
  id: z.number().optional().nullable().nullish(),
  name: z.string(),
  actions: z.string(),
  created_at: z.date().optional().nullable().nullish(),
  updated_at: z.date().optional().nullable().nullish(),
  deleted_at: z.date().optional().nullable().nullish(),
});

export const actionSchema = z.object({
  menu_dashboard: z.boolean().optional(),
  menu_sender: z.boolean().optional(),
  menu_config_user: z.boolean().optional(),
  menu_config_role: z.boolean().optional(),
  menu_system_flag: z.boolean().optional(),
  menu_system_log: z.boolean().optional(),
  menu_system_setting: z.boolean().optional(),
});

export const createRoleSchema = roleSchema
  .omit({
    id: true,
    actions: true,
    created_at: true,
    updated_at: true,
    deleted_at: true,
  })
  .extend({
    actions: actionSchema.optional().nullable(),
  });

export type Role = z.infer<typeof roleSchema>;
export type CreateRole = z.infer<typeof createRoleSchema>;
export type Action = z.infer<typeof actionSchema>;
