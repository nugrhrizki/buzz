import { createMutation, createQuery } from "@tanstack/solid-query";
import { Accessor } from "solid-js";

import { json, request } from "@/lib/request";
import { jidToPhoneNumber } from "@/lib/utils";

import { ContactResponse, Sender } from "@/models/sender";

const querySender = () => request.get<Array<Sender>>("/api/v1/whatsapp/users");
const querySenderDetail = (id: number) => request.get<Sender>(`/api/v1/whatsapp/user/${id}`);
const createQuerySender = (sender: Sender) => request.post<Sender>("/api/v1/whatsapp/create-user", sender);
const updateQuerySender = (sender: Sender) => request.put<Sender>(`/api/v1/whatsapp/update-user/${sender.id}`, sender);
const deleteQuerySender = (id: number) => request.delete<Sender>(`/api/v1/whatsapp/delete-user/${id}`);

export async function statusRequest(token: string) {
  if (token == "") {
    return null;
  }

  const response = await request.get("/api/v1/whatsapp/status", {
    headers: {
      token,
    },
  });
  return response.data;
}

export async function getQR(token: string) {
  if (token == "") {
    return null;
  }

  const response = await request.get("/api/v1/whatsapp/qr", {
    headers: {
      token,
    },
  });
  return response.data;
}

export async function connect(token: string) {
  if (token == "") {
    return null;
  }

  const response = await request.post(
    "/api/v1/whatsapp/connect",
    {
      events: "All",
      immediate: true,
    },
    {
      headers: {
        token,
      },
    },
  );
  return response.data;
}

export async function getContacts(token: string) {
  if (token == "") {
    return null;
  }

  const response = await request.post<ContactResponse>(
    "/api/v1/whatsapp/contacts",
    {},
    {
      headers: {
        token,
      },
    },
  );
  return response.data;
}

export function getAvatar(token: string, jid: string) {
  if (token == "") {
    throw json({
      title: "Oops, something went wrong!",
      message: "Token is required",
      statusCode: 400,
    });
  }

  const phone = jidToPhoneNumber(jid);
  if (!phone) {
    throw json({
      title: "Oops, something went wrong!",
      message: "Invalid phone number",
      statusCode: 400,
    });
  }

  return request.post(
    "/api/v1/whatsapp/avatar",
    { phone },
    {
      headers: {
        token,
      },
    },
  );
}

export function useSender() {
  return createQuery(() => ({
    queryKey: ["sender"],
    queryFn: async () => {
      const response = await querySender();
      return response.data;
    },
  }));
}

export function useSenderDetail(id: number) {
  return createQuery(() => ({
    queryKey: ["sender-detail", id],
    queryFn: async () => {
      const response = await querySenderDetail(id);
      return response.data;
    },
  }));
}

export function useAvatar(token: Accessor<string>, jid: Accessor<string>) {
  return createQuery(() => ({
    queryKey: ["sender-avatar", token(), jid()],
    queryFn: async () => {
      const response = await getAvatar(token(), jid());
      return response.data;
    },
    retry: 0,
    refetchOnMount: false,
    refetchInterval: false,
    refetchOnReconnect: false,
    refetchIntervalInBackground: false,
    refetchOnWindowFocus: false,
    gcTime: Infinity,
    networkMode: "offlineFirst",
    retryOnMount: false,
    staleTime: Infinity,
  }));
}

export function createSenderMutation() {
  return createMutation(() => ({
    mutationKey: ["sender"],
    mutationFn: async (sender: Sender) => {
      const response = await createQuerySender(sender);
      return response.data;
    },
  }));
}

export function updateSenderMutation() {
  return createMutation(() => ({
    mutationKey: ["sender"],
    mutationFn: async (sender: Sender) => {
      const response = await updateQuerySender(sender);
      return response.data;
    },
  }));
}

export function deleteSenderMutation() {
  return createMutation(() => ({
    mutationKey: ["sender"],
    mutationFn: async (id: number) => {
      const response = await deleteQuerySender(id);
      return response.data;
    },
  }));
}
