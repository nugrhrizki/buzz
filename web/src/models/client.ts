import { z } from "zod";

export const clientSchema = z.object({
  id: z.number().optional().nullable().nullish(),
  name: z.string(),
  token: z.string(),
  webhook: z.string().optional().nullable().nullish(),
  jid: z.string().optional().nullable().nullish(),
  qrcode: z.string().optional().nullable().nullish(),
  connected: z.number().optional().nullable().nullish(),
  expiration: z.number().optional().nullable().nullish(),
  events: z.string().optional().nullable().nullish(),
});

export const createClientSchema = clientSchema.omit({
  id: true,
  webhook: true,
  jid: true,
  qrcode: true,
  connected: true,
  expiration: true,
  events: true,
});

export type Client = z.infer<typeof clientSchema>;
export type CreateClient = z.infer<typeof createClientSchema>;
