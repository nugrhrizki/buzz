import { z } from "zod";

export const senderSchema = z.object({
  id: z.number().optional().nullable().nullish(),
  name: z.string().min(1, {
    message: "Name is required",
  }),
  token: z.string().min(1, {
    message: "Token is required",
  }),
  webhook: z.string().optional().nullable().nullish(),
  jid: z.string().optional().nullable().nullish(),
  qrcode: z.string().optional().nullable().nullish(),
  connected: z.number().optional().nullable().nullish(),
  expiration: z.number().optional().nullable().nullish(),
  events: z.string().optional().nullable().nullish(),
});

export const createSenderSchema = senderSchema.omit({
  id: true,
  webhook: true,
  jid: true,
  qrcode: true,
  connected: true,
  expiration: true,
  events: true,
});

export type Sender = z.infer<typeof senderSchema>;
export type CreateSender = z.infer<typeof createSenderSchema>;

export interface Contact {
  jid: string;
  found: boolean;
  firstName: string;
  fullName: string;
  pushName: string;
  businessName: string;
}

interface APIContact {
  Found: boolean;
  FirstName: string;
  FullName: string;
  PushName: string;
  BusinessName: string;
}

export type ContactResponse = Record<string, APIContact>;
